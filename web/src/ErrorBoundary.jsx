import { Button, HStack, Text, Textarea, useToast, VStack } from '@chakra-ui/react';
import React, { useState } from 'react';
import { GiBurningSkull } from 'react-icons/gi';
import { RiEmotionUnhappyLine } from 'react-icons/ri';
import Server from './server';
import { Center } from '@chakra-ui/react';


function ErrorToast(props) {
  const { error } = props
  const [showDetails, setShowDetails] = useState(false)

  const details = showDetails ?
    <VStack spacing={3}>
      <Text>{error.message}</Text>
      <Textarea rows={5}>{JSON.stringify(error.response)}</Textarea>
    </VStack> :
    <Button onClick={_ => setShowDetails(true)} colorScheme="blue">
      Show Technical Details
    </Button>

  return <VStack spacing={3} margin={5}>
    <Button onClick={_ => window.location.reload(false)} colorScheme="blue">
      Reload (sometimes it works)
    </Button>
    {details}
  </VStack>

}

let skullRef = null

function ServerError(props) {
  const toast = useToast()

  Server.addErrorHandler(0, r => {
    if (r.response.status > 499) {
      if (skullRef) {
        skullRef.style.display = 'block'
      }

      toast({
        title: <HStack>
          <Text>Something went wrong on the server</Text>
          <RiEmotionUnhappyLine />
        </HStack>,
        description: <ErrorToast error={r} />,
        status: "error",
        isClosable: true,
      })
      return null
    } else {
      return r
    }
  })
  return <Center style={{ display: 'none' }} ref={r => skullRef = r}>
    <GiBurningSkull size="30%" />
  </Center>
}

function ClientError(props) {
  const toast = useToast()
  const { errorInfo } = props

  if (errorInfo) {
    toast({
      title: <HStack>
        <Text>Something went wrong on the client</Text>
        <RiEmotionUnhappyLine />
      </HStack>,
      description: <ErrorToast error={errorInfo} />,
      status: "error",
      isClosable: true,
    })
  }

  return errorInfo ? <Center><GiBurningSkull size="30%" /></Center> : null;
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


class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { errorInfo: null };
    Server.addErrorHandler(10, loginWhenUnauthorized)
  }

  static getDerivedStateFromError(error) {
    return { hasError: true };
  }

  componentDidCatch(error, errorInfo) {
    this.setState({ errorInfo: errorInfo });
  }

  render() {
    if (this.state.hasError || this.state.errorInfo) {
      return <ClientError errorInfo={this.state.errorInfo} />;
    }

    return <>
      <ServerError />
      {this.props.children}
    </>;
  }
}
export default ErrorBoundary;