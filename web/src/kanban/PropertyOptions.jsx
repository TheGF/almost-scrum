import { Button, ButtonGroup, Checkbox, HStack, Menu, MenuButton, MenuItem, MenuList, Spacer, Stack, Text } from '@chakra-ui/react';
import Board from '@lourenci/react-kanban';
import '@lourenci/react-kanban/dist/styles.css';
import { React, useContext, useEffect, useState } from 'react';
import T from '../core/T';
import Server from '../server';
import UserContext from '../UserContext';
import MarkdownEditor from '../core/MarkdownEditor';



function PropertyOptions(props) {
    const { project } = useContext(UserContext)
    const { updateBoard, property, values } = props
    const [selectedValues, setSelectedValues] = useState([])


    function propertyOptionsDragEnd(board, card, source, destination) {
        if (board) {
            const t = board.columns[destination.toColumnId - 1]
            const ref = card.ref
            ref.task.properties[property] = t.title
            Server.setTask(project, ref.board, card.id, ref.task)
        }
    }

    function checkValue(v) {
        const values = selectedValues.includes(v) ?
            selectedValues.filter(s => s != v) :
            [...selectedValues, v]
        Server.postQueryTasks(project, {
            select: { description: true, properties: true },
            whereTypes: [{ hasPropertiesAll: [property] }]
        }).then(refs => {
            refs = refs || []
            updateBoard(values, refs, r => r.task.properties[property], propertyOptionsDragEnd)
        })
        setSelectedValues(values)

    }

    const valuesUI = values.map(v => <Checkbox key={v} onChange={_ => checkValue(v)}
        isChecked={selectedValues.includes(v)}>
        {v}
    </Checkbox>)

    return <HStack spacing={3}>{valuesUI}</HStack>
}

export default PropertyOptions