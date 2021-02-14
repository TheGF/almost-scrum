import {
    Button, Menu, MenuButton, MenuDivider, MenuItem,
    MenuList, Stack, useColorMode
} from "@chakra-ui/react";
import { React, useContext, useEffect, useState } from "react";
import { BiChevronDown } from "react-icons/all";
import T from "../core/T";
import UserContext from '../UserContext';
import Settings from './Settings';
import Help from '../help/Help';

const visibleBoards = 5

function Boards(props) {
    const { boards, active, onSelectBoard, onSelectLibrary, setSearchKeys } = props;
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
        const sorted = boards && boards.map(addTimeInfo)
            .sort((a, b) => b.tm - a.tm)
            .map(s => s.name) || []

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
        onSelectBoard && onSelectBoard("~")
    }

    function clickLibrary() {
        onSelectLibrary && onSelectLibrary()
    }

    const all = <Button key="all" colorScheme="blue"
        isActive={active == '~'} onClick={clickAll}>
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

    const library = <Button key="library" colorScheme="yellow"
        isActive={active == 'library'} onClick={clickLibrary}>
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
    const { info } = useContext(UserContext)

    const [activeBoard, setActiveBoard] = useState('backlog');
    const [boardKey, setBoardKey] = useState(0);
    const [showSettings, setShowSettings] = useState(false);
    const [showHelp, setShowHelp] = useState(false);

    const { onNewTask, onNewBoard, setShowGitIntegration, onExit } = props;
    const { colorMode, toggleColorMode } = useColorMode()

    function onSelectBoard(board) {
        setActiveBoard(board);
        props.setActiveBoard && props.setActiveBoard(board);
    }

    return <Stack spacing={4} direction="row" align="center">
        <Settings isOpen={showSettings} onClose={_ => setShowSettings(false)} />
        <Help isOpen={showHelp} onClose={_ => setShowHelp(false)} />
        <Menu>
            <MenuButton as={Button} rightIcon={<BiChevronDown />}>
                Actions
                </MenuButton>
            <MenuList>
                <MenuItem onClick={onNewBoard}>New Board</MenuItem>
                <MenuDivider />
                <MenuItem onClick={toggleColorMode}>
                    Toggle {colorMode === "light" ? "Dark" : "Light"}
                </MenuItem>
                <MenuDivider />
                {info && info.gitProject ? <MenuItem
                    onClick={_ => setShowGitIntegration(true)}>
                    Git Integration
                    </MenuItem> : null}
                <MenuItem onClick={_ => setShowSettings(true)}>
                    Users
                </MenuItem>
                <MenuItem onClick={_ => setShowHelp(true)}>
                    Help
                </MenuItem>
                {onExit ? <>
                    <MenuDivider />
                    <MenuItem onClick={onExit}><T>back to portal</T></MenuItem>
                </> : null}
            </MenuList>
        </Menu>
        <Boards key={boardKey}
            active={activeBoard} setActiveBoard={setActiveBoard}
            {...props} />
    </Stack>
}

export default Header;