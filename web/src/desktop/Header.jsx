import {
    Box, Button, Icon, IconButton, Menu, MenuButton, MenuDivider, MenuGroup, MenuItem,
    MenuItemOption, MenuList, MenuOptionGroup, Spacer, Stack, useColorMode, HStack, Text,
} from "@chakra-ui/react";
import { React, useContext, useState, useEffect } from "react";
import { BiChevronDown, BiUserVoice, BsChatSquareDots, GoTools, GrFormAdd } from "react-icons/all";
import T from "../core/T";
import Federation from "../federation/Federation";
import Help from '../help/Help';
import UserContext from '../UserContext';
import Boards from './Boards';
import Users from "./Users";
import Portal from '../portal/Portal';
import TodoList from './TodoList';
import TodoEdit from './TodoEdit';
import Server from "../server";
import Chat from '../chat/Chat';


function Header(props) {
    const { info, project, reload } = useContext(UserContext)
    const { colorMode, toggleColorMode } = useColorMode()

    const [boardKey, setBoardKey] = useState(0);
    const [showUsers, setShowUsers] = useState(false);
    const [showHelp, setShowHelp] = useState(false);
    const [editItem, setEditItem] = useState(null);
    const [todoList, setTodoList] = useState([])

    const lae = localStorage.getItem(`ash-options`)
    const [options, setOptions] = useState(
        lae && lae.split(',') || ['ui-effects', colorMode]
    );

    const { panel, setPanel, setShowGitIntegration, onExit, askBoardName } = props;

    function getTodoList() {
        Server.getUser(project, info.loginUser)
            .then(userInfo => {
                if (userInfo && userInfo.todo) {
                    const now = new Date();
                    const list = userInfo.todo.filter(t => !t.done || now < t.eta + 2*3600*1000)
                    setTodoList(list)
                }
            })
    }
    useEffect(getTodoList, [])

    function putTodoList(todoList) {
        Server.getUser(project, info.loginUser)
            .then(userInfo => userInfo && Server.setUser(project, info.loginUser, {...userInfo, todo: todoList}))

    }

    function optionsChange(value) {
        localStorage.setItem(`ash-options`, value.join(','))
        setOptions(value)
    }

    function confirmEditItem(editItem) {
        const list = todoList.map(t => t.id == editItem.id ? editItem : t)
        setTodoList(list)
        putTodoList(list)
        setEditItem(null)
    }


    function deleteEditItem(editItem) {
        const list = todoList.filter(t => t.id != editItem.id)
        setTodoList(list)
        putTodoList(list)
        setEditItem(null)
    }

    function changeTodoList(value, editItem) {
        setTodoList(value)
        if (editItem) {
            setEditItem(editItem)
        } else {
            putTodoList(value)    
        }
    }

    const library = <Button key="library" colorScheme="yellow"
        isActive={panel == '#library'} onClick={_ => setPanel('#library')}>
        <T>library</T>
    </Button>

    const ActionsMenu = <Menu>
        <MenuButton as={Button} colorScheme="blue" rightIcon={<BiChevronDown />}>
            Actions
        </MenuButton>
        <MenuList >
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
            <MenuItem onClick={_ => setPanel('#gantt')} command="⌘G">
                Gantt
            </MenuItem>
            <MenuItem onClick={_ => setPanel('#kanban')} command="⌘K">
                Kanban
            </MenuItem>
            <MenuDivider />
            <TodoList value={todoList} onChange={changeTodoList} />
        </MenuList>
    </Menu>

    return <Box id="header" w="90%">
        <Stack spacing={4} m={1} direction="row" align="center">
            <TodoEdit value={editItem} onConfirm={confirmEditItem} onDelete={deleteEditItem}/>
            <Users isOpen={showUsers} onClose={_ => setShowUsers(false)} />
            <Help isOpen={showHelp} onClose={_ => setShowHelp(false)} />
            {ActionsMenu}
            <Boards key={boardKey} panel={panel} setPanel={setPanel}
                {...props} />
            {library}
            {toolsMenu}
            <Spacer />
            <Chat/>
            <Federation />
        </Stack>
    </Box>
}

export default Header;