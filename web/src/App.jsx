import { Button, ChakraProvider, HStack, Text, useToast, VStack } from '@chakra-ui/react';
import { React, useEffect, useState } from 'react';
import theme from './theme'
import Desktop from './desktop/Desktop';
import Server from './server';
import Portal from './portal/Portal';
import ErrorBoundary from './ErrorBoundary'
import { RiEmotionUnhappyLine } from 'react-icons/ri';


function App() {
  const [portal, setPortal] = useState(null)
  const toast = useToast()

  function chooseMode() {
    Server.isPortal()
      .then(setPortal)
  }
  useEffect(chooseMode, [])

  const entry = portal == null ? null :
    portal ? <Portal /> : <Desktop project="~" />;
  return (
    <ChakraProvider theme={theme}  >
      <ErrorBoundary>{entry}</ErrorBoundary>
    </ChakraProvider>
  );
}

export default App;
