import { Box, Flex, useDisclosure, VStack } from '@chakra-ui/react';
import { React, useEffect, useState } from "react";
import Board from '../board/Board';
import Gantt from '../extensions/Gantt';
import GitIntegration from '../git/GitIntegration';
import Library from '../library/Library';
import Server from '../server';
import UserContext from '../UserContext';
import AskBoardName from './AskBoardName';
import Header from './Header';
import NoAccess from './NoAccess';
import Kanban from '../kanban/Kanban';
import ReactPlayer from 'react-player';


function Desktop(props) {
    const { project, onExit } = props;

    const [info, setInfo] = useState(null)
    const [panel, setPanel] = useState(null)
    const [boards, setBoards] = useState([]);
    const [showGitIntegration, setShowGitIntegration] = useState(false);
    const [noAccess, setNoAccess] = useState(null)
    const askBoardName = useDisclosure(false)
    const [refreshId, setRefreshId] = useState(false)

    function checkNoAccess(r) {
        if (r.response && r.response.status == 403 && r.response.data.users && r.response.data.message) {
            setNoAccess(r.response.data)
            return r
        }
        return r
    }

    function exit() {
        setNoAccess(null)
        onExit()
    }

    function listBoards() {
        Server.listBoards(project)
            .then(setBoards)
    }

    function createBoard(name) {
        Server.createBoard(project, name)
            .then(askBoardName.onClose())
            .then(listBoards)
            .then(setPanel(name))
    }


    function init() {
        Server.addErrorHandler(10, checkNoAccess)
        Server.getProjectInfo(project)
            .then(info => {
                setInfo(info)
                listBoards()
            })
            .catch(checkNoAccess)
    }
    useEffect(init, [])

    function getContent() {
        switch (panel) {
            case null:
            case "": return null;
            case '#library':
                return <Library />
            case '#gantt':
                return <Gantt />
            case '#kanban':
                return <Kanban />
            default:
                return <Board key={panel} name={panel} boards={boards} />
        }
    }


    const reload = () => window.location.reload();
    const username = info && info.systemUser
    const userContext = { project, info, username, reload }

    const desktop = <div key={refreshId} >
        <AskBoardName {...askBoardName} boards={boards} onCreate={createBoard} />
        <GitIntegration isOpen={showGitIntegration}
            onClose={_ => setShowGitIntegration(false)} />

        <Flex
            direction="column"
            align="center"
            w="100%"
            spacing={10}
            m="1 auto">

            <Header boards={boards} setShowGitIntegration={setShowGitIntegration}
                panel={panel} setPanel={setPanel}
                onListBoards={listBoards} askBoardName={askBoardName}
                onExit={onExit} />
            <Box w="90%" mt={5}>
                {getContent()}
            </Box>


            {/* <Center p="1em" style={{ position: "fixed", bottom: 0 }}>
            <Text size="sm" color="gray">
                {info && `${info.loginUser}@${project}.${info.host}`}
             - Scrum To Go 0.5 - Open Source since 2020</Text>
        </Center> */}
        </Flex>
    </div>
    const body = info ? <UserContext.Provider value={userContext}>
        {desktop}
    </UserContext.Provider> : null

    return noAccess ? <NoAccess data={noAccess} onExit={exit} /> : body

}

export default Desktop