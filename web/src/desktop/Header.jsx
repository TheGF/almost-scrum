import {
    Button, Menu, MenuButton, MenuDivider, MenuItem,
    MenuList, Stack, useColorMode
} from "@chakra-ui/react";
import { React, useContext, useState } from "react";
import { BiChevronDown } from "react-icons/all";
import T from "../core/T";
import Help from '../help/Help';
import UserContext from '../UserContext';
import Boards from './Boards';
import Settings from './Settings';


function Header(props) {
    const { info } = useContext(UserContext)

    const [activeBoard, setActiveBoard] = useState(null);
    const [boardKey, setBoardKey] = useState(0);
    const [showSettings, setShowSettings] = useState(false);
    const [showHelp, setShowHelp] = useState(false);

    const { setShowGitIntegration, onExit, askBoardName } = props;
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
                <MenuItem onClick={_ => askBoardName.onOpen()}>
                    New Board
                </MenuItem>
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