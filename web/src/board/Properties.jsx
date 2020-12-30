import { Input, Table, Tbody, Td, Tr, VStack } from '@chakra-ui/react';
import { React, useState } from "react";
import T from '../core/T';

function Properties(props) {
    const { task, setTask } = props;

    function renderProperty(key, val) {
        const [value, setValue] = useState(val);

        function onChange(evt) {
            const value = evt && evt.target && evt.target.value;
            if (value == undefined) return

            task.properties[key] = value
            setValue(value)
            setTask(task)
        }

        return <Tr key={key}>
            <Td><T>{key}</T></Td>
            <Td>
                <Input value={value} onChange={onChange}/>
            </Td>
        </Tr>
    }

    const rows = task.properties && Object.entries(task.properties).map(
        e => renderProperty(e[0], e[1])
    )

    return <Table variant="striped" colorScheme="teal" size="sm">
        <Tbody>
            {rows}
        </Tbody>
    </Table>;
}
export default Properties;