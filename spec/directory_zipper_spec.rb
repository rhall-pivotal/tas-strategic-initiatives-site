require 'spec_helper'
require 'directory_zipper'

describe DirectoryZipper do
  around do |test|
    FileUtils.rm_rf(temp_dir)
    FileUtils.cp_r(File.join(fixture_dir, '.'), temp_dir)
    test.run
    FileUtils.rm_rf(temp_dir)
  end

  let(:common_source_directory) { File.join(temp_dir, 'common_source_directory') }
  let(:nested_directory_with_contents) { File.join(common_source_directory, 'nested_directory_with_contents') }
  let(:sibling_directory_to_nested_directory) { File.join(common_source_directory, 'sibling_directory_to_nested_directory') }
  let(:directory_in_another_world) { File.join(temp_dir, 'directory_in_another_world') }
  let(:empty_directory) { File.join(temp_dir, 'empty_directory') }

  let(:target_zip) { File.join(temp_dir, 'directory_zipper_testcase.zip') }

  let(:directory_zipper) { DirectoryZipper.new(target_zip, common_source_directory) }

  describe "#zip" do
    it "creates a zip file at the target zip" do
      expect {
        directory_zipper.zip
      }.to change { File.exists?(target_zip) }.from(false).to(true)
    end

    context 'when a directory and a file are added' do
      let(:expected_zip_structure) {
        %w(
            nested_directory_with_contents/
            nested_directory_with_contents/nested_file.one
            sibling_directory_to_nested_directory/
            sibling_directory_to_nested_directory/sibling_file.one
            sibling_directory_to_nested_directory/sibling_file.two
          )
      }

      it 'zips up both the directory and the file' do
        directory_zipper.add_file(File.join(nested_directory_with_contents, 'nested_file.one'))
        directory_zipper.add_directory(sibling_directory_to_nested_directory)
        directory_zipper.zip

        expect(contents_of_zipfile(target_zip)).to match_array(expected_zip_structure)
      end
    end
  end

  describe "#add_directory" do
    describe "when adding the common source directory to the zip" do
      let(:common_source_directory) { directory_in_another_world }
      let(:expected_zip_structure) {
        %w(
          far_far_away.file
        )
      }

      it "should not include the common source directory as an entry (only the contents)" do
        directory_zipper.add_directory(common_source_directory)
        directory_zipper.zip
        expect(contents_of_zipfile(target_zip)).to match_array(expected_zip_structure)
      end
    end

    describe "when the directory shares the common source directory" do
      let(:expected_zip_structure) {
        %w(
          top_level_file.one
          top_level_file.two
          nested_directory_with_contents/
          nested_directory_with_contents/nested_file.one
          nested_directory_with_contents/nested_file.two
          sibling_directory_to_nested_directory/
          sibling_directory_to_nested_directory/sibling_file.one
          sibling_directory_to_nested_directory/sibling_file.two
        )
      }

      it "recursively adds directory and all its contents to the zip preserving any nested structures" do
        directory_zipper.add_directory(common_source_directory)
        directory_zipper.zip
        expect(contents_of_zipfile(target_zip)).to match_array(expected_zip_structure)
      end

      describe "adding an empty directory" do
        let(:common_source_directory) { temp_dir }
        let(:expected_zip_structure) {
          %w(
            empty_directory/
          )
        }

        it "should zip up an empty directory" do
          directory_zipper.add_directory(empty_directory)
          directory_zipper.zip
          expect(contents_of_zipfile(target_zip)).to match_array(expected_zip_structure)
        end
      end

      describe "adding multiple directories" do
        let(:expected_zip_structure) {
          %w(
            nested_directory_with_contents/
            nested_directory_with_contents/nested_file.one
            nested_directory_with_contents/nested_file.two
            sibling_directory_to_nested_directory/
            sibling_directory_to_nested_directory/sibling_file.one
            sibling_directory_to_nested_directory/sibling_file.two
          )
        }

        it "recursively adds the new directory and all its contents" do
          directory_zipper.add_directory(nested_directory_with_contents)
          directory_zipper.add_directory(sibling_directory_to_nested_directory)
          directory_zipper.zip
          expect(contents_of_zipfile(target_zip)).to match_array(expected_zip_structure)
        end
      end
    end

    describe "adding a directory that does not share the common source directory" do
      it "should raise" do
        expect {
          directory_zipper.add_directory(directory_in_another_world)
        }.to raise_error(/Must add files and directories under the common source directory/)
      end
    end
  end

  describe "#add_file" do
    let(:file_in_nested_directory) { File.join(nested_directory_with_contents, 'nested_file.one') }
    let(:file_in_sibling_directory) { File.join(sibling_directory_to_nested_directory, 'sibling_file.one') }
    let(:file_in_a_different_common_directory) { File.join(directory_in_another_world, 'far_far_away.file')}

    describe "when the file shares the common_source_directory" do
      let(:expected_zip_structure) {
        %w(
          nested_directory_with_contents/
          nested_directory_with_contents/nested_file.one
        )
      }

      it "should be added to the zip with the directory structure past the common_source_directory" do
        directory_zipper.add_file(file_in_nested_directory)
        directory_zipper.zip
        expect(contents_of_zipfile(target_zip)).to match_array(expected_zip_structure)
      end
    end

    describe "when adding multiple files" do
      let(:expected_zip_structure) {
        %w(
          nested_directory_with_contents/
          nested_directory_with_contents/nested_file.one
          sibling_directory_to_nested_directory/
          sibling_directory_to_nested_directory/sibling_file.one
        )
      }

      it "should zip both files" do
        directory_zipper.add_file(file_in_nested_directory)
        directory_zipper.add_file(file_in_sibling_directory)
        directory_zipper.zip
        expect(contents_of_zipfile(target_zip)).to match_array(expected_zip_structure)
      end
    end

    describe "when the file does not share the common_source_directory" do
      it "should raise" do
        expect {
          directory_zipper.add_file(file_in_a_different_common_directory)
        }.to raise_error(/Must add files and directories under the common source directory/)
      end
    end
  end
end
