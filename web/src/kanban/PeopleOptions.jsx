import { Button, ButtonGroup, Checkbox, HStack, Menu, MenuButton, MenuItem, MenuList, Spacer, Stack, Text } from '@chakra-ui/react';
import Board from '@lourenci/react-kanban';
import '@lourenci/react-kanban/dist/styles.css';
import { React, useContext, useEffect, useState } from 'react';
import T from '../core/T';
import Server from '../server';
import UserContext from '../UserContext';
import MarkdownEditor from '../core/MarkdownEditor';



function PeopleOptions(props) {
    const { project } = useContext(UserContext)
    const { updateBoard } = props
    const [users, setUsers] = useState([])
    const [selectedUsers, setSelectedUsers] = useState([])

    function listUsers() {
        Server.listUsers(project)
            .then(setUsers)
    }
    useEffect(_ => listUsers(), [])


    function peopleOptionsDragEnd(board, card, source, destination) {
        if (board) {
            const t = board.columns[destination.toColumnId - 1]
            const ref = card.ref
            ref.task.properties['Owner'] = `@${t.title}`
            Server.setTask(project, ref.board, card.id, ref.task)
        }
    }

    function checkValue(v) {
        const values = selectedUsers.includes(v) ?
            selectedUsers.filter(s => s != v) :
            [...selectedUsers, v]
        Server.postQueryTasks(project, {
            select: { description: true, properties: true },
            whereTypes: [{ hasPropertiesAll: ['Owner'] }]
        }).then(refs => {
            refs = refs || []
            updateBoard(values, refs, r => r.task.properties['Owner'].replace(/^@/, ''), peopleOptionsDragEnd)
        })
        setSelectedUsers(values)

    }

    const usersUI = users.map(v => <Checkbox key={v} onChange={_ => checkValue(v)}
        isChecked={selectedUsers.includes(v)}>
        {v}
    </Checkbox>)

    return <HStack spacing={3}>{usersUI}</HStack>
}

export default PeopleOptions