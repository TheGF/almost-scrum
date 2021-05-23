import { Button, ButtonGroup, Checkbox, HStack, Menu, MenuButton, MenuItem, MenuList, Spacer, Stack, Text } from '@chakra-ui/react';
import Board from '@lourenci/react-kanban';
import '@lourenci/react-kanban/dist/styles.css';
import { React, useContext, useEffect, useState } from 'react';
import T from '../core/T';
import Server from '../server';
import UserContext from '../UserContext';
import MarkdownEditor from '../core/MarkdownEditor';
import BoardOptions from './BoardOptions';
import PropertyOptions from './PropertyOptions';
import PeopleOptions from './PeopleOptions';

function Options(props) {
    const { info } = useContext(UserContext)
    const initialViewId = localStorage.getItem('kanban-view') || '!boards'
    const { updateBoard } = props
    const [viewId, setViewId] = useState(initialViewId)
    const properties = getProperties()

    function getProperties() {
        const vp = {}
        for (const m of info.models) {
            for (const p of m.properties) {
                if (['Tag', 'Enum'].includes(p.kind) && p.values) {
                    vp[p.name] = p.values
                }
            }
        }
        return vp
    }

    function getViewSelection() {
        const propertiesUI = Object.keys(properties).map(p => <MenuItem key={p} onClick={_=>updateViewId(p)}>
            {p}
        </MenuItem>)
        return <HStack>
            <ButtonGroup size="sm" >
                <Button isActive={viewId == '!boards'} onClick={_ => updateViewId('!boards')}><T>boards</T></Button>
                <Button isActive={viewId == '!people'} onClick={_ => updateViewId('!people')}><T>people</T></Button>
                <Menu>
                    <MenuButton as={Button} >
                        Property
                    </MenuButton>
                    <MenuList>
                        {propertiesUI}
                    </MenuList>
                </Menu>
            </ButtonGroup>
        </HStack>
    }

    function updateViewId(viewId) {
        localStorage.setItem('kanban-view', viewId)
        setViewId(viewId)
    }

    function getViewById(viewId) {
        switch (viewId) {
            case '!boards': return <BoardOptions updateBoard={updateBoard} />
            case '!people': return <PeopleOptions updateBoard={updateBoard} />
            default:
                return <PropertyOptions updateBoard={updateBoard} property={viewId} values={properties[viewId]} />
        }
    }

    return <>
        {getViewSelection()}
        {getViewById(viewId)}
    </>

}

export default Options