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
    const { boards, panel, setPanel, onListBoards } = props;
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

    function clickBoard(panel) {
        localStorage.setItem(`ash-board-${panel}`, `${Date.now()}`)
        setPanel(panel)
        splitBoards(recentBoards)
    }

    function clickAll() {
        setActiveButton('all')
        onSelectBoard && onSelectBoard("~")
    }

    function closeEditBoard() {
        setEditBoard(null)
    }

    const all = <Button key="~" 
        isActive={panel == '~'} onClick={_=>clickBoard('~')}>
        <T>all</T>
    </Button>

    const buttons = recentBoards.map(
        b => <Button key={b} isActive={b == panel}
            onClick={_ => clickBoard(b)} >
            <T>{b}</T>
            <Box w="1em" />
            <Link onClick={_ => setEditBoard(b)}><BiEdit /></Link>
        </Button>
    );

    const moreButtons = more.length ?
        <Menu>
            <MenuButton as={Button}>
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


    return <>
        <EditBoard name={editBoard} onClose={closeEditBoard} onListBoards={onListBoards} />
        {buttons}
        {moreButtons}
        {all}
    </>
}
export default Boards