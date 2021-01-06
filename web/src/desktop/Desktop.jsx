import { Flex, useDisclosure, VStack } from '@chakra-ui/react';
import { React, useContext, useEffect, useState } from "react";
import Board from '../board/Board';
import Library from '../library/Library';
import Server from '../server';
import UserContext from '../UserContext';
import AskBoardName from './AskBoardName';
import Header from './Header';

function Desktop() {
    const { project, username } = useContext(UserContext);
    const [board, setBoard] = useState('backlog');
    const [boards, setBoards] = useState([]);
    const [boardKey, setBoardKey] = useState(0);
    const [showLibrary, setShowLibrary] = useState(false);
    const askBoardName = useDisclosure(false)

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
        Server.createTask(project, board, 'Click_and_Rename', {
            description: "",
            properties: {
                "Owner": `@${username}`,
                "Status": "Draft",
            },
            progress: [],
            attachments: [],
        })
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

    return <Flex
        direction="column"
        align="center"
        w={{ xl: "83%" }}
        m="0 auto">

        <AskBoardName {...askBoardName} boards={boards} onCreate={createBoard} />
        <VStack w="100%">
            <Header boards={boards}
                onSelectBoard={onSelectBoard} onSelectLibrary={onSelectLibrary}
                onNewTask={onNewTask} onNewBoard={_ => askBoardName.onOpen()} />
            {content}
        </VStack>
    </Flex>
}

export default Desktop