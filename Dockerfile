FROM docker.vsphere.gocd.cf-app.com:5000/releng:releng_base

ADD include/id_rsa /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa

ADD Gemfile /Gemfile
ADD Gemfile.lock /Gemfile.lock

RUN ssh-agent bundle

# We have dependency conflicts with bosh_cli
RUN gem install bosh_cli

#Don't publish ssh key
RUN rm /root/.ssh/id_rsa
