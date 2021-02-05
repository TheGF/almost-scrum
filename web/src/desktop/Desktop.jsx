import { Flex, useDisclosure, VStack } from '@chakra-ui/react';
import { React, useContext, useEffect, useState } from "react";
import Board from '../board/Board';
import Library from '../library/Library';
import Server from '../server';
import UserContext from '../UserContext';
import AskBoardName from './AskBoardName';
import Header from './Header';
import GitIntegration from '../git/GitIntegration';


function Desktop(props) {
    const { project, onExit } = props;
    const [info, setInfo] = useState(null)
    const [board, setBoard] = useState('backlog');
    const [boards, setBoards] = useState([]);
    const [boardKey, setBoardKey] = useState(0);
    const [showLibrary, setShowLibrary] = useState(false);
    const [showGitIntegration, setShowGitIntegration] = useState(false);

    const askBoardName = useDisclosure(false)

    function getInfo() {
        Server.getProjectInfo(project)
            .then(setInfo)
    }
    useEffect(getInfo, [])

    function listBoards() {
        Server.listBoards(project)
            .then(setBoards)
    }
    useEffect(listBoards, []);

    function createBoard(name) {
        Server.createBoard(project, name)
            .then(askBoardName.onClose())
            .then(listBoards)
            .then(setBoard(name))
    }

    function onNewTask() {
        Server.createTask(project, board, 'Click_and_Rename')
            .then(_ => setBoardKey(1 + boardKey))
    }

    function onSelectLibrary() {
        setBoard(null);
        setShowLibrary(true);
    }

    function onSelectBoard(board) {
        setBoard(board);
        setShowLibrary(false);
    }

    const content = showLibrary ? <Library /> :
        <Board key={boardKey} name={board} boards={boards} />

    const username  = info && info.systemUser
    const userContext = { project, info, username }
    return info ? <UserContext.Provider value={userContext}>

        <GitIntegration isOpen={showGitIntegration}
            onClose={_ => setShowGitIntegration(false)} />

        <AskBoardName {...askBoardName} boards={boards} onCreate={createBoard} />
        <Flex
            direction="column"
            align="center"
            w={{ xl: "83%" }}
            m="0 auto">

            <VStack w="100%">
                <Header boards={boards}
                    setShowGitIntegration={setShowGitIntegration}
                    onSelectBoard={onSelectBoard} onSelectLibrary={onSelectLibrary}
                    onNewTask={onNewTask} onNewBoard={_ => askBoardName.onOpen()}
                    onExit={onExit} />
                {content}
            </VStack>
        </Flex>
    </UserContext.Provider> : null
}

export default Desktop