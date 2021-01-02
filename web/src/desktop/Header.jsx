import {
    Button, Container, Menu, MenuButton, MenuItem, MenuList,
    MenuOptionGroup, MenuItemOption, MenuDivider, Stack, Select, MenuGroup
} from "@chakra-ui/react";
import { React, useEffect, useState, useContext } from "react";
import { BiChevronDown } from "react-icons/all";
import UserContext from '../UserContext'
import Server from "../server";
import T from "../core/T"

const visibleBoards = 5

function Boards(props) {
    const { boards, active, onSelectBoard, onSelectLibrary } = props;
    const [recentBoards, setRecentBoards] = useState([]);
    const [more, setMore] = useState([]);

    function splitBoards() {
        function addTimeInfo(b) {
            if (b == 'backlog') {
                return {
                    name: b,
                    tm: Date.now()
                }
            } else {
                return {
                    name: b,
                    tm: parseInt(localStorage.getItem(`ash-board-${b}`) || '0'),
                }
            }
        }
        const sorted = boards.map(addTimeInfo)
            .sort((a, b) => b.tm - a.tm)
            .map(s => s.name)

        setRecentBoards(sorted.slice(0, 5))
        setMore(sorted.slice(5))
    }
    useEffect(splitBoards, [boards])

    function clickBoard(board) {
        localStorage.setItem(`ash-board-${board}`, `${Date.now()}`)
        onSelectBoard && onSelectBoard(board)
        splitBoards(recentBoards)
    }

    function clickAll() {
        onSelectBoard && onSelectBoard("")
    }

    function clickLibrary() {
        onSelectLibrary && onSelectLibrary()
    }

    const all = <Button key="all" colorScheme="blue"
        isActive={active == ''} onClick={clickAll}>
        <T>all</T>
    </Button>

    const buttons = recentBoards.map(
        b => <Button key={b} colorScheme="blue" isActive={active == b}
            onClick={_ => clickBoard(b)} >
            <T>{b}</T>
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

    const library = <Button key="library" colorScheme="yellow" isActive={active == 'library'} >
        <T>library</T>
    </Button>

    return <>
        {buttons}
        {moreButtons}
        {all}
        {library}
    </>
}

function Header(props) {
    const [activeBoard, setActiveBoard] = useState('backlog');
    const [boardKey, setBoardKey] = useState(0);
    const {onNewTask, onNewBoard} = props;

    function onSelectBoard(board) {
        setActiveBoard(board);
        props.setActiveBoard && props.setActiveBoard(board);
    }

    return <Stack spacing={4} direction="row" align="center">
        <Menu>
            <MenuButton as={Button} rightIcon={<BiChevronDown />}>
                Actions
                </MenuButton>
            <MenuList>
                <MenuItem onClick={onNewTask}>New Task</MenuItem>
                <MenuItem onClick={onNewBoard}>New Board</MenuItem>
                <MenuDivider />
                <MenuOptionGroup title="Look & Feel" type="checkbox">
                    <MenuItemOption value="asc">Dark Mode</MenuItemOption>
                </MenuOptionGroup>
                <MenuDivider />
                <MenuItem>Git Push</MenuItem>
                <MenuItem>Git Pull</MenuItem>
            </MenuList>
        </Menu>
        <Boards key={boardKey} {...props}
            active={activeBoard} setActiveBoard={setActiveBoard} />
    </Stack>
}

export default Header;