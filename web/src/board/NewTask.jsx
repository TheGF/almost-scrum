import { Button, HStack, Spacer, IconButton, Menu, MenuButton, MenuList, MenuItem } from "@chakra-ui/react";
import {
    Popover,
    PopoverTrigger,
    PopoverContent,
    PopoverHeader,
    PopoverBody,
    PopoverFooter,
    PopoverArrow,
    PopoverCloseButton,
} from "@chakra-ui/react"
import { React, useContext, useEffect, useRef, useState } from "react";
import { BsViewStacked, MdViewHeadline, RiChatNewLine, RiFilterLine } from 'react-icons/all';
import ReactTags from 'react-tag-autocomplete';
import './reactTags.css'
import Server from '../server';
import UserContext from '../UserContext';
import Filter from "./Filter";
import Portal from '../portal/Portal';

function NewTask(props) {
    const { info } = useContext(UserContext);
    const { models, config } = info || {}
    const { board, onNewTask } = props

    const types = (config.boardTypes && config.boardTypes[board]) || models.map(m => m.name)
    const typesList = types.map(t => <MenuItem key={t} onClick={_ => onNewTask(t)}>
        {t}
    </MenuItem>)

    return types.length == 1 ?
        <IconButton title="New Task" icon={<RiChatNewLine />}
            onClick={_ => onNewTask(types[0])} /> :
        <Menu>
            <MenuButton as={Button} title="New Task">
                <RiChatNewLine />
            </MenuButton>
            <MenuList>
                {typesList}
            </MenuList>
        </Menu>
}

export default NewTask
