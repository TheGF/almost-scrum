import {
    Box,
    Button, IconButton, Link, Menu, MenuButton, MenuItem,
    MenuList
} from "@chakra-ui/react";
import { React, useEffect, useState } from "react";
import { BiEdit } from "react-icons/bi";
import T from "../core/T";
import EditBoard from './EditBoard';

const visibleBoards = 5

function Boards(props) {
    const { activeBoard, boards, onSelectBoard, onSelectLibrary, onListBoards } = props;
    const [activeButton, setActiveButton] = useState(activeBoard)
    const [recentBoards, setRecentBoards] = useState([]);
    const [more, setMore] = useState([]);
    const [editBoard, setEditBoard] = useState(null);

    function splitBoards() {
        function addTimeInfo(b) {
            return {
                name: b,
                tm: parseInt(localStorage.getItem(`ash-board-${b}`) || '0'),
            }
        }
        const sorted = boards && boards.map(addTimeInfo)
            .sort((a, b) => b.tm - a.tm)
            .map(s => s.name) || []

        setRecentBoards(sorted.slice(0, visibleBoards))
        setMore(sorted.slice(visibleBoards))
    }
    useEffect(splitBoards, [boards])

    function clickBoard(board) {
        setActiveButton(board)
        localStorage.setItem(`ash-board-${board}`, `${Date.now()}`)
        onSelectBoard && onSelectBoard(board)
        splitBoards(recentBoards)
    }

    function clickAll() {
        setActiveButton('all')
        onSelectBoard && onSelectBoard("~")
    }

    function clickLibrary() {
        setActiveButton('library')
        onSelectLibrary && onSelectLibrary()
    }

    function closeEditBoard() {
        setEditBoard(null)
    }

    const all = <Button key="all" colorScheme="blue"
        isActive={activeButton == 'library'} onClick={clickAll}>
        <T>all</T>
    </Button>

    const buttons = recentBoards.map(
        b => <Button key={b} colorScheme="blue" isActive={b == activeButton}
            onClick={_ => clickBoard(b)} >
            <T>{b}</T>
            <Box w="1em" />
            <Link onClick={_ => setEditBoard(b)}><BiEdit /></Link>
        </Button>
    );

    const moreButtons = more.length ?
        <Menu>
            <MenuButton as={Button} colorScheme="blue">
                ...
            </MenuButton>
            <MenuList>
                {
                    more.map(
                        b => <MenuItem key={b} onClick={_ => clickBoard(b)}>
                            {b}
                        </MenuItem>
                    )
                }
            </MenuList>
        </Menu> : null;

    const library = <Button key="library" colorScheme="yellow"
        isActive={activeButton == 'library'} onClick={clickLibrary}>
        <T>library</T>
    </Button>

    return <>
        <EditBoard name={editBoard} onClose={ closeEditBoard } onListBoards={onListBoards}/>
        {buttons}
        {moreButtons}
        {all}
        {library}
    </>
}
export default Boards