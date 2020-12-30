import { ChakraProvider, theme } from '@chakra-ui/react';
import { React } from 'react';
import Desktop from './desktop/Desktop';
import UserContext from './UserContext';


function App() {
  const project = '~'
  const value = { project }
  return (
    <UserContext.Provider value={value}>
      <ChakraProvider theme={theme}>
        <Desktop />
      </ChakraProvider>
    </UserContext.Provider>
  );
}

export default App;
