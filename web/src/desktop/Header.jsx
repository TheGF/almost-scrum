import {
    Button, Menu, MenuButton, MenuDivider, MenuItem,
    MenuList, Stack, useColorMode
} from "@chakra-ui/react";
import { React, useContext, useState } from "react";
import { BiChevronDown } from "react-icons/all";
import T from "../core/T";
import Federation from "../federation/Federation";
import Help from '../help/Help';
import UserContext from '../UserContext';
import Boards from './Boards';
import Users from "./Users";


function Header(props) {
    const { info } = useContext(UserContext)

    const [selectedPanel, setSelectedPanel] = useState(null);
    const [boardKey, setBoardKey] = useState(0);
    const [showUsers, setShowUsers] = useState(false);
    const [showHelp, setShowHelp] = useState(false);

    const { setShowGitIntegration, onExit, askBoardName } = props;
    const { colorMode, toggleColorMode } = useColorMode()

    function selectPanel(panel) {
        setSelectedPanel(panel);
        props.selectPanel && props.selectPanel(panel);
    }

    const gantt = <Button key="gantt" colorScheme="yellow"
        isActive={selectedPanel == '#gantt'} onClick={_ => selectPanel('#gantt')} >
        <T>gantt</T>
    </Button>

    const library = <Button key="library" colorScheme="yellow"
        isActive={selectedPanel == '#library'} onClick={_ => selectPanel('#library')}>
        <T>library</T>
    </Button>


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
                <Users />
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
            active={selectedPanel} setActiveBoard={setSelectedPanel}
            {...props} />
        {gantt}
        {library}
        <Federation />
    </Stack>
}

export default Header;