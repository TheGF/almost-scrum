import { Center, Flex, Spacer, Text, useDisclosure, VStack } from '@chakra-ui/react';
import { React, useEffect, useState } from "react";
import Board from '../board/Board';
import Gantt from '../gantt/Gantt';
import GitIntegration from '../git/GitIntegration';
import Library from '../library/Library';
import Server from '../server';
import UserContext from '../UserContext';
import AskBoardName from './AskBoardName';
import Header from './Header';
import NoAccess from './NoAccess';


function Desktop(props) {
    const { project, onExit } = props;

    const [info, setInfo] = useState(null)
    const [content, setContent] = useState(null)
    const [activeBoard, setActiveBoard] = useState(null)
    const [boards, setBoards] = useState([]);
    const [showLibrary, setShowLibrary] = useState(false);
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
            .then(setActiveBoard(name))
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

    function onSelectLibrary() {
        setShowLibrary(true);
    }

    function selectPanel(panel) {
        switch(panel) {
            case '#library': 
                setContent(<Library/>); break
            case '#gantt':
                setContent(<Gantt/>); break
            default:
                setContent(<Board name={panel} boards={boards} />); break
            }        
    }

    const reload = () => setRefreshId(!refreshId);
    const username = info && info.systemUser
    const userContext = { project, info, username, reload }
    const body = info ? <UserContext.Provider value={userContext}>
        <div key={refreshId}>
            <AskBoardName {...askBoardName} boards={boards} onCreate={createBoard} />
            <GitIntegration isOpen={showGitIntegration}
                onClose={_ => setShowGitIntegration(false)} />

            <Flex
                direction="column"
                align="center"
                w={{ xl: "83%" }}
                m="0 auto">

                <VStack w="100%" h="100%">
                    <Header boards={boards} setShowGitIntegration={setShowGitIntegration}
                        selectPanel={selectPanel}
                        onListBoards={listBoards} askBoardName={askBoardName}
                        onExit={onExit} />
                    {content}
                </VStack>
                <Center p="1em" style={{ position: "fixed", bottom: 0 }}>
                    <Text size="sm" color="gray">
                     {info && `${info.loginUser}@${project}.${info.host}`} 
                     - Almost Scrum 0.5 - Open Source since 2020</Text>
                </Center>
            </Flex>
        </div>
    </UserContext.Provider> : null

    return noAccess ? <NoAccess data={noAccess} onExit={exit} /> : body

}

export default Desktop