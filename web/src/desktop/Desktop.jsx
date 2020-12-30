import { React, useEffect, useState, useContext } from "react";
import { Flex, VStack } from '@chakra-ui/react';
import Header from './Header'
import Board from '../board/Board'
import FilterPanel from './FilterPanel';

function Desktop() {

    const [board, setBoard] = useState('backlog');
    const [showLibrary, setShowLibrary] = useState(false);

    function onSelectLibrary() {
        setBoard(null);
        setShowLibrary(true);
    }

    function onSelectBoard(board) {
        setBoard(board);
        setShowLibrary(false);
    }

    return <Flex
        direction="column"
        align="center"
        maxW={{ xl: "1200px" }}
        m="0 auto">
        <VStack>
            <Header onSelectBoard={onSelectBoard} onSelectLibrary={onSelectLibrary} />
            <Board name={board} />
        </VStack>
    </Flex>
}

export default Desktop