import {
  Button, Modal, ModalBody, ModalContent,
  ModalFooter, ModalHeader, ModalOverlay
} from '@chakra-ui/react';
import { React } from 'react';


function ConfirmUpload(props) {
  const { file, uploadFileToLibrary, getNextVersion } = props;

  function upgradeVersionAndUpload() {
    const [prefix, version, ext] = getNextVersion(file)
    uploadFileToLibrary(file, `${prefix}${version}${ext}`)
  }

  return <Modal isOpen={file}>
    <ModalOverlay />
    <ModalContent>
      <ModalHeader>File Exists</ModalHeader>
      <ModalBody>
        A version of {file && file.name} is already in the library.
      </ModalBody>

      <ModalFooter>
        <Button colorScheme="blue" mr={3} onClick={_ => uploadFileToLibrary(file)}>
          Overwrite
        </Button>
        <Button colorScheme="blue" mr={3} onClick={upgradeVersionAndUpload}>
          New Version
        </Button>
        <Button mr={3} onClick={_ => uploadFileToLibrary(null)}>
          Cancel
        </Button>
      </ModalFooter>
    </ModalContent>
  </Modal>
}

export default ConfirmUpload