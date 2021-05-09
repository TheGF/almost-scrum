import {
    Button, Center, HStack, Input, Spacer, Switch, Table,
    TableCaption, Tbody, Td, Tr
} from '@chakra-ui/react';
import { React, useState } from "react";
import { FiPlusCircle, VscTrash } from 'react-icons/all';

function Progress(props) {
    const { task, saveTask, readOnly } = props;
    const [parts, setParts] = useState(task.parts || []);

    function RenderPart(props) {
        const { part } = props;
        const [description, setDescription] = useState(part.description)
        const [done, setDone] = useState(part.done || false);

        function onDescriptionChange(evt) {
            const value = evt && evt.target && evt.target.value;
            if (value == undefined) return

            part.description = value
            setDescription(value)
        }

        function onDoneChange(evt) {
            part.done = !done
            setDone(part.done)
            saveTask(task)
        }

        function addRow() {
            const idx = task.parts.indexOf(part)
            task.parts = [
                ...parts.slice(0, 1 + idx),
                {
                    description: '',
                    done: false,
                },
                ...parts.slice(1 + idx)
            ]
            setParts(task.parts)
            saveTask(task)
        }


        return <Tr>
            <Td w="2em">
                <HStack spacing="1">
                    <Button onClick={addRow} disabled={readOnly}><FiPlusCircle /></Button>
                    <Button ><VscTrash /></Button>
                </HStack>
            </Td>
            <Td>
                <Input value={description} readOnly={readOnly}
                    placeholder="Describe a required action"
                    onChange={onDescriptionChange} onBlur={_=>saveTask(task)} />
            </Td>
            <Td w="2em">
                <Switch isChecked={done} isReadOnly={readOnly}
                    onChange={onDoneChange} />
            </Td>
        </Tr>
    }

    function addFirstRow() {
        task.parts = [{
            description: '',
            done: false,
        }]
        setParts(task.parts)
        saveTask(task)
    }

    const rows = parts.map((part, idx) => <RenderPart key={idx} part={part} />)
    const add = parts.length ? null : <Button onClick={addFirstRow}><FiPlusCircle /></Button>
    const editMessage = readOnly ? <Center h="2em">
        Change owner if you want to edit the content
    </Center> : null

    return <div className="panel2" >
        {editMessage}
        {add}
        <Table size="sm">
            <TableCaption w="100%">
                <HStack spacing={3}>
                    <Center>Track your progress and acceptance criteria</Center>
                    <Spacer />
                </HStack>
            </TableCaption>
            <Tbody>
                {rows}
            </Tbody>
        </Table>
    </div>;
}
export default Progress;