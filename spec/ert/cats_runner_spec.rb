require 'spec_helper'
require 'ert/cats_runner'

describe Ert::CatsRunner do
  let(:environment_name) { 'env_name' }
  let(:env_config) do
    <<YAML
---
name: env_name
iaas_type: vsphere
YAML
  end

  let(:settings) do
    RecursiveOpenStruct.new(YAML.load(env_config), recurse_over_arrays: true)
  end

  let(:opsmgr_environment) { instance_double(Opsmgr::Environments, settings: settings) }
  let(:bosh_command) { instance_double(Opsmgr::Cmd::BoshCommand) }
  let(:deployment_name) { 'cf-deadbeef12345678' }
  let(:logger) { instance_double(Opsmgr::LoggerWithProgName) }

  subject(:cats_runner) { Ert::CatsRunner.new(environment_name: environment_name, om_version: '1.5', logger: logger) }

  before do
    allow(Opsmgr::Environments).to receive(:for).and_return(opsmgr_environment)
    allow(Opsmgr::Cmd::BoshCommand).to(
      receive(:new)
        .and_return(bosh_command)
    )
    allow(logger).to receive(:info)
  end

  it 'uses the config for the given environment' do
    expect(Opsmgr::Environments).to receive(:for).with(environment_name)
    cats_runner
  end

  context '#run_cats' do
    before do
      allow(bosh_command).to receive(:target).and_return('the bosh target command')
      allow(bosh_command).to receive(:command).and_return('bosh_command')
      allow(ENV).to receive(:fetch).with('TMPDIR', '/tmp').and_return('temp_dir')
      allow(ENV).to receive(:[]=).and_call_original
      allow(Open3).to receive(:capture2).and_return(["#{deployment_name}\n", instance_double(Process::Status, success?: true)])
      allow(Bundler).to receive(:clean_system).and_return(true)
    end

    context 'vsphere' do
      it 'directly targets the microbosh' do
        expect(bosh_command).to receive(:target).and_return('the bosh target command')
        expect(Bundler).to receive(:clean_system).with('the bosh target command')

        cats_runner.run_cats
      end

      it 'downloads the manifest to TMPDIR/env_name.yml' do
        expect(ENV).to receive(:fetch).with('TMPDIR', '/tmp').and_return('temp_dir')
        expect(Open3).to(
          receive(:capture2)
            .with("bosh_command deployments | grep -Eoh 'cf-[0-9a-f]{8,}'")
            .and_return(["#{deployment_name}\n", instance_double(Process::Status, success?: true)])
        )
        expect(Bundler).to(
          receive(:clean_system)
            .with("bosh_command -n download manifest #{deployment_name} temp_dir/#{environment_name}.yml")
        )

        cats_runner.run_cats
      end

      it 'sets the deployment to the downloaded manifest' do
        expect(Bundler).to receive(:clean_system).with("bosh_command deployment temp_dir/#{environment_name}.yml")

        cats_runner.run_cats
      end

      context 'when running the cats errand' do
        context 'in an environment with full internet connectivity' do
          it 'runs the acceptance-tests errand' do
            expect(Bundler).to receive(:clean_system).with('bosh_command run errand acceptance-tests')

            cats_runner.run_cats
          end

          context 'when the acceptance-tests errand fails' do
            it 'raises a RuntimeError' do
              expect(Bundler).to receive(:clean_system).with('bosh_command run errand acceptance-tests').and_return(false)

              expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'CF Acceptance Tests failed')
            end
          end
        end

        context 'in an environment with no internet connectivity' do
          let(:env_config) do
            <<YAML
---
name: env_name
iaas_type: vsphere
internetless: true
YAML
          end

          it 'runs the acceptance-tests errand' do
            expect(Bundler).to receive(:clean_system).with('bosh_command run errand acceptance-tests-internetless')

            cats_runner.run_cats
          end

          context 'when the acceptance-tests-internetless errand fails' do
            it 'raises a RuntimeError' do
              expect(Bundler).to(
                receive(:clean_system).with('bosh_command run errand acceptance-tests-internetless')
                  .and_return(false)
              )
              expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'CF Acceptance Tests failed')
            end
          end
        end
      end

      context 'when a bosh command fails' do
        context 'when bosh target fails' do
          it 'raises a RuntimeError' do
            expect(Bundler).to receive(:clean_system).with('the bosh target command').and_return(false)

            expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'bosh target failed')
          end
        end

        context 'when bosh deployments fails' do
          it 'raises a RuntimeError' do
            expect(Open3).to(
              receive(:capture2)
                .and_return(["#{deployment_name}\n", instance_double(Process::Status, success?: false)])
            )

            expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'bosh deployments failed')
          end
        end

        context 'when bosh download manifest fails' do
          it 'raises a RuntimeError' do
            expect(Bundler).to(
              receive(:clean_system)
                .with("bosh_command -n download manifest #{deployment_name} temp_dir/#{environment_name}.yml")
            ).and_return(false)

            expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'bosh download manifest failed')
          end
        end

        context 'when bosh deployment fails' do
          it 'raises a RuntimeError' do
            expect(Bundler).to(
              receive(:clean_system).with("bosh_command deployment temp_dir/#{environment_name}.yml")
                .and_return(false)
            )

            expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'bosh deployment failed')
          end
        end
      end
    end

    context 'aws' do
      let(:env_config) do
        <<YAML
---
name: env_name
iaas_type: aws
ops_manager:
  url: https://foo.com
  aws:
    ssh_key: key
YAML
      end

      let(:gateway) { instance_double(Net::SSH::Gateway) }

      before do
        allow(Net::SSH::Gateway).to receive(:new).and_return(gateway)
        allow(gateway).to receive(:open).and_yield(25_555)
        allow(bosh_command).to receive(:director_ip).and_return('1.2.3.4')
      end

      it 'sets up the ssh gateway' do
        expect(Net::SSH::Gateway).to receive(:new).with('foo.com', 'ubuntu', key_data: ['key']).and_return(gateway)
        expect(gateway).to receive(:open).with('1.2.3.4', 25_555, 25_555).and_yield(25_555)
        expect(bosh_command).to receive(:director_ip).and_return('1.2.3.4')

        cats_runner.run_cats
      end

      it 'sets the DIRECTOR_IP_OVERRIDE environment var' do
        expect(ENV).to receive(:[]=).with('DIRECTOR_IP_OVERRIDE', 'localhost')
        cats_runner.run_cats
      end

      it 'targets the microbosh' do
        expect(bosh_command).to receive(:target).and_return('the bosh target command')
        expect(Bundler).to receive(:clean_system).with('the bosh target command')

        cats_runner.run_cats
      end

      it 'downloads the manifest to TMPDIR/env_name.yml' do
        expect(ENV).to receive(:fetch).with('TMPDIR', '/tmp').and_return('temp_dir')
        expect(Open3).to(
          receive(:capture2)
            .with("bosh_command deployments | grep -Eoh 'cf-[0-9a-f]{8,}'")
            .and_return(["#{deployment_name}\n", instance_double(Process::Status, success?: true)])
        )
        expect(Bundler).to(
          receive(:clean_system)
            .with("bosh_command -n download manifest #{deployment_name} temp_dir/#{environment_name}.yml")
        )

        cats_runner.run_cats
      end

      it 'sets the deployment to the downloaded manifest' do
        expect(Bundler).to receive(:clean_system).with("bosh_command deployment temp_dir/#{environment_name}.yml")

        cats_runner.run_cats
      end

      context 'when running the cats errand' do
        context 'in an environment with full internet connectivity' do
          it 'runs the acceptance-tests errand' do
            expect(Bundler).to receive(:clean_system).with('bosh_command run errand acceptance-tests')

            cats_runner.run_cats
          end

          context 'when the acceptance-tests errand fails' do
            it 'raises a RuntimeError' do
              expect(Bundler).to receive(:clean_system).with('bosh_command run errand acceptance-tests').and_return(false)

              expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'CF Acceptance Tests failed')
            end
          end
        end

        context 'in an environment with no internet connectivity' do
          let(:env_config) do
            <<YAML
---
name: env_name
iaas_type: vsphere
internetless: true
YAML
          end

          it 'runs the acceptance-tests errand' do
            expect(Bundler).to receive(:clean_system).with('bosh_command run errand acceptance-tests-internetless')

            cats_runner.run_cats
          end

          context 'when the acceptance-tests-internetless errand fails' do
            it 'raises a RuntimeError' do
              expect(Bundler).to(
                receive(:clean_system).with('bosh_command run errand acceptance-tests-internetless')
                  .and_return(false)
              )
              expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'CF Acceptance Tests failed')
            end
          end
        end
      end

      context 'when a bosh command fails' do
        context 'when bosh target fails' do
          it 'raises a RuntimeError' do
            expect(Bundler).to receive(:clean_system).with('the bosh target command').and_return(false)

            expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'bosh target failed')
          end
        end

        context 'when bosh deployments fails' do
          it 'raises a RuntimeError' do
            expect(Open3).to(
              receive(:capture2)
                .and_return(["#{deployment_name}\n", instance_double(Process::Status, success?: false)])
            )

            expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'bosh deployments failed')
          end
        end

        context 'when bosh download manifest fails' do
          it 'raises a RuntimeError' do
            expect(Bundler).to(
              receive(:clean_system)
                .with("bosh_command -n download manifest #{deployment_name} temp_dir/#{environment_name}.yml")
            ).and_return(false)

            expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'bosh download manifest failed')
          end
        end

        context 'when bosh deployment fails' do
          it 'raises a RuntimeError' do
            expect(Bundler).to(
              receive(:clean_system).with("bosh_command deployment temp_dir/#{environment_name}.yml")
                .and_return(false)
            )

            expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'bosh deployment failed')
          end
        end
      end
    end

    context 'vcloud' do
      let(:env_config) do
        <<YAML
---
name: env_name
iaas_type: vcloud
ops_manager:
  url: https://foo.com
YAML
      end

      let(:gateway) { instance_double(Net::SSH::Gateway) }

      before do
        allow(Net::SSH::Gateway).to receive(:new).and_return(gateway)
        allow(gateway).to receive(:open).and_yield(25_555)
        allow(bosh_command).to receive(:director_ip).and_return('1.2.3.4')
      end

      it 'sets up the ssh gateway' do
        expect(Net::SSH::Gateway).to receive(:new).with('foo.com', 'ubuntu', password: 'tempest').and_return(gateway)
        expect(gateway).to receive(:open).with('1.2.3.4', 25_555, 25_555).and_yield(25_555)
        expect(bosh_command).to receive(:director_ip).and_return('1.2.3.4')

        cats_runner.run_cats
      end

      it 'sets the DIRECTOR_IP_OVERRIDE environment var' do
        expect(ENV).to receive(:[]=).with('DIRECTOR_IP_OVERRIDE', 'localhost')
        cats_runner.run_cats
      end

      it 'targets the microbosh' do
        expect(bosh_command).to receive(:target).and_return('the bosh target command')
        expect(Bundler).to receive(:clean_system).with('the bosh target command')

        cats_runner.run_cats
      end

      it 'downloads the manifest to TMPDIR/env_name.yml' do
        expect(ENV).to receive(:fetch).with('TMPDIR', '/tmp').and_return('temp_dir')
        expect(Open3).to(
          receive(:capture2)
            .with("bosh_command deployments | grep -Eoh 'cf-[0-9a-f]{8,}'")
            .and_return(["#{deployment_name}\n", instance_double(Process::Status, success?: true)])
        )
        expect(Bundler).to(
          receive(:clean_system)
            .with("bosh_command -n download manifest #{deployment_name} temp_dir/#{environment_name}.yml")
        )

        cats_runner.run_cats
      end

      it 'sets the deployment to the downloaded manifest' do
        expect(Bundler).to receive(:clean_system).with("bosh_command deployment temp_dir/#{environment_name}.yml")

        cats_runner.run_cats
      end

      context 'when running the cats errand' do
        context 'in an environment with full internet connectivity' do
          it 'runs the acceptance-tests errand' do
            expect(Bundler).to receive(:clean_system).with('bosh_command run errand acceptance-tests')

            cats_runner.run_cats
          end

          context 'when the acceptance-tests errand fails' do
            it 'raises a RuntimeError' do
              expect(Bundler).to receive(:clean_system).with('bosh_command run errand acceptance-tests').and_return(false)

              expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'CF Acceptance Tests failed')
            end
          end
        end

        context 'in an environment with no internet connectivity' do
          let(:env_config) do
            <<YAML
---
name: env_name
iaas_type: vsphere
internetless: true
YAML
          end

          it 'runs the acceptance-tests errand' do
            expect(Bundler).to receive(:clean_system).with('bosh_command run errand acceptance-tests-internetless')

            cats_runner.run_cats
          end

          context 'when the acceptance-tests-internetless errand fails' do
            it 'raises a RuntimeError' do
              expect(Bundler).to(
                receive(:clean_system).with('bosh_command run errand acceptance-tests-internetless')
                  .and_return(false)
              )
              expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'CF Acceptance Tests failed')
            end
          end
        end
      end

      context 'when a bosh command fails' do
        context 'when bosh target fails' do
          it 'raises a RuntimeError' do
            expect(Bundler).to receive(:clean_system).with('the bosh target command').and_return(false)

            expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'bosh target failed')
          end
        end

        context 'when bosh deployments fails' do
          it 'raises a RuntimeError' do
            expect(Open3).to(
              receive(:capture2)
                .and_return(["#{deployment_name}\n", instance_double(Process::Status, success?: false)])
            )

            expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'bosh deployments failed')
          end
        end

        context 'when bosh download manifest fails' do
          it 'raises a RuntimeError' do
            expect(Bundler).to(
              receive(:clean_system)
                .with("bosh_command -n download manifest #{deployment_name} temp_dir/#{environment_name}.yml")
            ).and_return(false)

            expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'bosh download manifest failed')
          end
        end

        context 'when bosh deployment fails' do
          it 'raises a RuntimeError' do
            expect(Bundler).to(
              receive(:clean_system).with("bosh_command deployment temp_dir/#{environment_name}.yml")
                .and_return(false)
            )

            expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'bosh deployment failed')
          end
        end
      end
    end

    context 'openstack' do
      let(:env_config) do
        <<YAML
---
name: env_name
iaas_type: openstack
ops_manager:
  url: https://foo.com
  openstack:
    ssh_private_key: key
YAML
      end

      let(:gateway) { instance_double(Net::SSH::Gateway) }

      before do
        allow(Net::SSH::Gateway).to receive(:new).and_return(gateway)
        allow(gateway).to receive(:open).and_yield(25_555)
        allow(bosh_command).to receive(:director_ip).and_return('1.2.3.4')
      end

      it 'sets up the ssh gateway' do
        expect(Net::SSH::Gateway).to receive(:new).with('foo.com', 'ubuntu', key_data: ['key']).and_return(gateway)
        expect(gateway).to receive(:open).with('1.2.3.4', 25_555, 25_555).and_yield(25_555)
        expect(bosh_command).to receive(:director_ip).and_return('1.2.3.4')

        cats_runner.run_cats
      end

      it 'sets the DIRECTOR_IP_OVERRIDE environment var' do
        expect(ENV).to receive(:[]=).with('DIRECTOR_IP_OVERRIDE', 'localhost')
        cats_runner.run_cats
      end

      it 'targets the microbosh' do
        expect(bosh_command).to receive(:target).and_return('the bosh target command')
        expect(Bundler).to receive(:clean_system).with('the bosh target command')

        cats_runner.run_cats
      end

      it 'downloads the manifest to TMPDIR/env_name.yml' do
        expect(ENV).to receive(:fetch).with('TMPDIR', '/tmp').and_return('temp_dir')
        expect(Open3).to(
          receive(:capture2)
            .with("bosh_command deployments | grep -Eoh 'cf-[0-9a-f]{8,}'")
            .and_return(["#{deployment_name}\n", instance_double(Process::Status, success?: true)])
        )
        expect(Bundler).to(
          receive(:clean_system)
            .with("bosh_command -n download manifest #{deployment_name} temp_dir/#{environment_name}.yml")
        )

        cats_runner.run_cats
      end

      it 'sets the deployment to the downloaded manifest' do
        expect(Bundler).to receive(:clean_system).with("bosh_command deployment temp_dir/#{environment_name}.yml")

        cats_runner.run_cats
      end

      context 'when running the cats errand' do
        context 'in an environment with full internet connectivity' do
          it 'runs the acceptance-tests errand' do
            expect(Bundler).to receive(:clean_system).with('bosh_command run errand acceptance-tests')

            cats_runner.run_cats
          end

          context 'when the acceptance-tests errand fails' do
            it 'raises a RuntimeError' do
              expect(Bundler).to receive(:clean_system).with('bosh_command run errand acceptance-tests').and_return(false)

              expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'CF Acceptance Tests failed')
            end
          end
        end

        context 'in an environment with no internet connectivity' do
          let(:env_config) do
            <<YAML
---
name: env_name
iaas_type: vsphere
internetless: true
YAML
          end

          it 'runs the acceptance-tests errand' do
            expect(Bundler).to receive(:clean_system).with('bosh_command run errand acceptance-tests-internetless')

            cats_runner.run_cats
          end

          context 'when the acceptance-tests-internetless errand fails' do
            it 'raises a RuntimeError' do
              expect(Bundler).to(
                receive(:clean_system).with('bosh_command run errand acceptance-tests-internetless')
                  .and_return(false)
              )
              expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'CF Acceptance Tests failed')
            end
          end
        end
      end

      context 'when a bosh command fails' do
        context 'when bosh target fails' do
          it 'raises a RuntimeError' do
            expect(Bundler).to receive(:clean_system).with('the bosh target command').and_return(false)

            expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'bosh target failed')
          end
        end

        context 'when bosh deployments fails' do
          it 'raises a RuntimeError' do
            expect(Open3).to(
              receive(:capture2)
                .and_return(["#{deployment_name}\n", instance_double(Process::Status, success?: false)])
            )

            expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'bosh deployments failed')
          end
        end

        context 'when bosh download manifest fails' do
          it 'raises a RuntimeError' do
            expect(Bundler).to(
              receive(:clean_system)
                .with("bosh_command -n download manifest #{deployment_name} temp_dir/#{environment_name}.yml")
            ).and_return(false)

            expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'bosh download manifest failed')
          end
        end

        context 'when bosh deployment fails' do
          it 'raises a RuntimeError' do
            expect(Bundler).to(
              receive(:clean_system).with("bosh_command deployment temp_dir/#{environment_name}.yml")
                .and_return(false)
            )

            expect { cats_runner.run_cats }.to raise_error(RuntimeError, 'bosh deployment failed')
          end
        end
      end
    end
  end
end
