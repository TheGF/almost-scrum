import {
    Button, IconButton, Menu, MenuButton, MenuDivider, MenuItem,
    MenuList, Stack, useColorMode
} from "@chakra-ui/react";
import { React, useContext, useState } from "react";
import { AiOutlineCloudSync, BiChevronDown, BiTransfer, GiMeshBall, GiMeshNetwork, GiSpiderWeb } from "react-icons/all";
import T from "../core/T";
import Federation from "../federation/Federation";
import Help from '../help/Help';
import UserContext from '../UserContext';
import Boards from './Boards';
import Settings from './Settings';
import Users from "./Users";


function Header(props) {
    const { info } = useContext(UserContext)

    const [activeBoard, setActiveBoard] = useState(null);
    const [boardKey, setBoardKey] = useState(0);
    const [showUsers, setShowUsers] = useState(false);
    const [showHelp, setShowHelp] = useState(false);

    const { setShowGitIntegration, onExit, askBoardName } = props;
    const { colorMode, toggleColorMode } = useColorMode()

    function onSelectBoard(board) {
        setActiveBoard(board);
        props.setActiveBoard && props.setActiveBoard(board);
    }

    return <Stack spacing={4} direction="row" align="center">
        <Users isOpen={showUsers} onClose={_ => setShowUsers(false)} />
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
                <Users/>
                <MenuItem onClick={_ => setShowUsers(true)}>
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
        <Federation/>
    </Stack>
}

export default Header;