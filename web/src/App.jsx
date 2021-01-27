import { ChakraProvider } from '@chakra-ui/react';
import { React, useEffect, useState } from 'react';
import theme from './theme'
import Desktop from './desktop/Desktop';
import Server from './server';
import Portal from './portal/Portal';


function App() {
  const [portal, setPortal] = useState(null)

  function chooseMode() {
    Server.isPortal()
      .then(setPortal)
  }
  useEffect(chooseMode, [])

  const entry = portal == null ? null :
    portal ? <Portal /> : <Desktop project="~" />;
  return (
    <ChakraProvider theme={theme}  >
      {/* <Global styles={globalStyles} />
      <CSSReset /> */}
      {entry}
    </ChakraProvider>
  );
}

export default App;
