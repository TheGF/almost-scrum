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
    const [boards, setBoards] = useState([]);
    const [more, setMore] = useState([]);
    const [active, setActive] = useState('backlog');

    function splitBoards(boards) {
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
        boards = boards.map(addTimeInfo)
            .sort((a, b) => b.tm - a.tm)
            .map(s => s.name)

        setBoards(boards.slice(0, 5))
        setMore(boards.slice(5))
    }

    function clickBoard(board) {
        setActive(board)
        localStorage.setItem(`ash-board-${board}`, `${Date.now()}`)
        if (props.onSelectBoard) {
            props.onSelectBoard(board)
        }
        splitBoards(boards)
    }

    function clickAll() {
        setActive('all')
        if (props.onSelectBoard) {
            props.onSelectBoard('')
        }
    }

    function clickLibrary() {
        setActive('library')
        if (props.onSelectLibrary) {
            props.onSelectLibrary()
        }
    }

    function listBoards() {
        Server.listBoards("~")
            .then(splitBoards)
    }
    useEffect(listBoards, []);

    const all = <Button key="all" colorScheme="blue"
        isActive={active == 'all'} onClick={clickAll}>
        <T>all</T>
    </Button>

    const buttons = boards.map(
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
    const { project } = useContext(UserContext);

    return <Stack spacing={4} direction="row" align="center">
        <Menu>
            <MenuButton as={Button} rightIcon={<BiChevronDown />}>
                Actions
                </MenuButton>
            <MenuList>
                <MenuItem>New Task</MenuItem>
                <MenuItem>New Board</MenuItem>
                <MenuDivider />
                <MenuOptionGroup title="Look & Feel" type="checkbox">
                    <MenuItemOption value="asc">Dark Mode</MenuItemOption>
                </MenuOptionGroup>
                <MenuDivider />
                <MenuItem>Git Push</MenuItem>
                <MenuItem>Git Pull</MenuItem>
            </MenuList>
        </Menu>
        <Boards {...props} />
    </Stack>
}

export default Header;