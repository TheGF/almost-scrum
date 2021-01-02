import { ChakraProvider, theme } from '@chakra-ui/react';
import { React, useEffect, useState } from 'react';
import Desktop from './desktop/Desktop';
import UserContext from './UserContext';
import Server from './server';


function App() {
  const project = '~'

  const [info, setInfo] = useState(null)
  const username  = info && info.system_user
  const value = { project, info, username }

  function getInfo() {
    Server.getProjectInfo(project)
      .then(setInfo)
  }
  useEffect(getInfo, [])

  return (
    <UserContext.Provider value={value}>
      <ChakraProvider theme={theme}>
        <Desktop />
      </ChakraProvider>
    </UserContext.Provider>
  );
}

export default App;
