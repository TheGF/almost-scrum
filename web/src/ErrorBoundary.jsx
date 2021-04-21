import {
  Button, HStack, Modal, ModalBody, ModalContent, ModalHeader,
  Spacer,
  Text, Textarea, useToast, VStack
} from '@chakra-ui/react';
import React, { useState } from 'react';
import { GiBurningSkull } from 'react-icons/gi';
import { RiEmotionUnhappyLine } from 'react-icons/ri';
import Server from './server';
import { Center } from '@chakra-ui/react';


function ErrorToast(props) {
  const { errorInfo } = props
  const [showDetails, setShowDetails] = useState(false)

  const details = showDetails ?
    <VStack spacing={3} align="left">
      <Text><b>Message: </b>{errorInfo.message} {' '}
        {errorInfo.response && errorInfo.response.statusText || ''}</Text>
      <Text><b>Data: </b>{errorInfo.response && errorInfo.response.data || ''}</Text>
    </VStack> :
    <Button onClick={_ => setShowDetails(true)} colorScheme="blue">
      Show Details
    </Button>

  return errorInfo ? <HStack>
    <GiBurningSkull size="30%" />
    <Spacer/>
    <GiBurningSkull size="30%" />
    <Modal isOpen >
      <ModalContent>
        <ModalHeader>
          <HStack>
            <RiEmotionUnhappyLine />
            <Text>Something went wrong</Text>
          </HStack>
        </ModalHeader>
        <ModalBody>
          <VStack spacing={3} margin={5}>
            <Button onClick={_ => window.location.reload(false)} colorScheme="blue">
              Reload (sometimes it works)
            </Button>
            {details}
          </VStack>
        </ModalBody>
      </ModalContent>
    </Modal>
  </HStack> : null
}

function loginWhenUnauthorized(r) {
  if (r && r.response && r.response.status == 401) {
    localStorage.removeItem('username');
    localStorage.removeItem('token');
    window.location.assign(window.location.href);
  } else {
    return r
  }
}

function ShowMessage(props) {
  const { title, errorInfo } = props
  const toast = useToast()

  toast({
    title: <HStack>
      <Text>{title}</Text>
      <RiEmotionUnhappyLine />
    </HStack>,
    description: <ErrorToast error={errorInfo} />,
    status: "error",
    isClosable: true,
  })

  return
}

let errorBoundary = null

class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);

    this.state = { errorInfo: null };
    this.serverError = this.serverError.bind(this)
    this.errorOccured = false

    errorBoundary = this
    Server.addErrorHandler(0, r => errorBoundary.serverError(r));
  }

  serverError(r) {
    if (r && r.response && r.response.status == 401) {
      localStorage.removeItem('username');
      localStorage.removeItem('token');
      window.location.assign(window.location.href);
      return null
    }
    if (r && r.response && r.response.status > 499 && !this.errorOccured) {
      this.errorOccured = true
      this.setState({ hasError: true, errorInfo: r })
      return null
    } else {
      return r
    }
  }

  static getDerivedStateFromError(error) {
    return { hasError: true, errorInfo: error };
  }

  componentDidCatch(error, more) {
    this.setState({ errorInfo: {...error, data: more.componentStack }});
  }

  render() {
    return this.state.hasError ?
      <ErrorToast errorInfo={this.state.errorInfo} /> :
      this.props.children
  }
}
export default ErrorBoundary;