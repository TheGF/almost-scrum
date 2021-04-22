import {
    Box, Button, Icon, IconButton, Menu, MenuButton, MenuDivider, MenuItem,
    MenuItemOption, MenuList, MenuOptionGroup, Spacer, Stack, useColorMode
} from "@chakra-ui/react";
import { React, useContext, useState } from "react";
import { BiChevronDown, GoTools } from "react-icons/all";
import T from "../core/T";
import Federation from "../federation/Federation";
import Help from '../help/Help';
import UserContext from '../UserContext';
import Boards from './Boards';
import Users from "./Users";


function Header(props) {
    const { info, project, reload } = useContext(UserContext)
    const { colorMode, toggleColorMode } = useColorMode()

    const [boardKey, setBoardKey] = useState(0);
    const [showUsers, setShowUsers] = useState(false);
    const [showHelp, setShowHelp] = useState(false);

    const lae = localStorage.getItem(`ash-options`)
    const [options, setOptions] = useState(
        lae && lae.split(',') || ['ui-effects', colorMode]
    );

    const { panel, setPanel, setShowGitIntegration, onExit, askBoardName } = props;

    function optionsChange(value) {
        localStorage.setItem(`ash-options`, value.join(','))
        setOptions(value)
    }

    const library = <Button key="library" colorScheme="yellow"
        isActive={panel == '#library'} onClick={_ => setPanel('#library')}>
        <T>library</T>
    </Button>

    const ActionsMenu = <Menu>
        <MenuButton as={Button} colorScheme="blue" rightIcon={<BiChevronDown />}>
            Actions
        </MenuButton>
        <MenuList>
            <MenuItem onClick={_ => askBoardName.onOpen()}>
                New Board
        </MenuItem>
            <MenuDivider />
            <MenuOptionGroup title="Options" type="checkbox"
                value={options} onChange={optionsChange}>
                <MenuItemOption value="ui-effects" onClick={reload}>
                    UI Effects
                </MenuItemOption>
                <MenuItemOption value="dark" onClick={toggleColorMode}>
                    Dark Mode
                </MenuItemOption>
            </MenuOptionGroup>

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

    const toolsMenu = <Menu>
        <MenuButton as={Button} colorScheme="yellow">
            <GoTools />
        </MenuButton>
        <MenuList>
            <MenuItem onClick={_ => setPanel('#gantt')}>
                Gantt
            </MenuItem>
            <MenuItem onClick={_ => setPanel('#kanban')}>
                Kanban
            </MenuItem>
        </MenuList>
    </Menu>

    return <Box id="header" w="90%">
        <Stack spacing={4} m={1} direction="row" align="center">
            <Users isOpen={showUsers} onClose={_ => setShowUsers(false)} />
            <Help isOpen={showHelp} onClose={_ => setShowHelp(false)} />
            {ActionsMenu}
            <Boards key={boardKey} panel={panel} setPanel={setPanel}
                {...props} />
            {library}
            {toolsMenu}
            <Spacer />
            <Federation />
        </Stack>
    </Box>
}

export default Header;