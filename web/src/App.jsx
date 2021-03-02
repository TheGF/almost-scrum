import { ChakraProvider } from '@chakra-ui/react';
import { React, useEffect, useState } from 'react';
import Desktop from './desktop/Desktop';
import ErrorBoundary from './ErrorBoundary';
import Portal from './portal/Portal';
import Server from './server';
import theme from './theme';

function uuidv4() {
  return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, c =>
    (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
  );
}
let clientId = uuidv4()

function App() {
  const [hello, setHello] = useState(null)

  function chooseMode() {
    Server.hello(clientId)
      .then(setHello)
      .then(_ => {
        window.addEventListener('beforeunload', function(){
          Server.bye(clientId);
          alert('BYE')
        })
      })
  }
  useEffect(chooseMode, [])

  const systemUser = hello && hello.systemUser
  const entry = hello == null ? null :
    hello.portal ? <Portal systemUser={systemUser} /> : <Desktop project="~" />;
  return (
    <ChakraProvider theme={theme}  >
      <ErrorBoundary>{entry}</ErrorBoundary>
    </ChakraProvider>
  );
}

export default App;
