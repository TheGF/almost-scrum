import {
    Center, Input, Select, Switch, Table, TableCaption, Tbody,
    Td, Tr, Box,
} from '@chakra-ui/react';
import { React, useContext, useState } from "react";
import T from '../core/T';
import UserContext from '../UserContext';

function Properties(props) {
    const { task, saveTask, readOnly, users, height } = props;
    const { info } = useContext(UserContext)
    const { propertyModel } = info
    const { properties } = task;

    function renderProperty(propertyDef) {
        const { name, kind, values } = propertyDef
        const [value, setValue] = useState(properties[name] || '');

        function onChange(evt) {
            const value = evt && evt.target && evt.target.value;
            if (value == undefined) return

            properties[name] = value
            setValue(value)
            saveTask(task)
        }

        function renderString() {
            return <Input value={value} onChange={onChange} size="sm"/>
        }

        function renderEnum() {
            const options = values && values.map(option =>
                <option value={option} key={option}>{T.translate(option)}</option>
            ) || []
            return readOnly ? <label>{value}</label> :
                <Select placeholder="Choose" value={value}
                    size="small" onChange={onChange}>
                    {options}
                </Select>
        }

        function renderTag() {
            const options = values && values.map(option =>
                <option value={option} key={option}>{option}</option>
            ) || []
            return readOnly ? <label>{value}</label> :
                <Select placeholder="Choose" value={value}
                    size="small" onChange={onChange}>
                    {options}
                </Select>
        }

        function renderUser() {
            const options = users && users.map(user =>
                <option value={`@${user}`} key={user}>{user}</option>
            ) || []
            const label = value && value.substring(1)
            return readOnly ? <label>{label}</label> :
                <Select placeholder="Choose" value={value}
                    size="small" onChange={onChange}>
                    {options}
                </Select>
        }

        function renderBool() {
            return readOnly ? <label>{value}</label> :
                <Switch isChecked={value} isReadOnly={readOnly}
                    onChange={onChange} />
        }

        let input = null

        switch (kind) {
            case 'String': input = renderString(); break
            case 'Enum': input = renderEnum(); break
            case 'Bool': input = renderBool(); break
            case 'User': input = renderUser(); break
            case 'Tag': input = renderTag(); break
        }

        return <Tr key={name}>
            <Td><T>{name}</T></Td>
            <Td>{input}</Td>
        </Tr>
    }

    const rows = (propertyModel || []).map(propertyDef => renderProperty(propertyDef))
    return <Box maxH={height-50} style={{ overflowY: 'auto' }}>
        <Table variant="striped" colorScheme="teal" size="sm">
            <Tbody>
                {rows}
            </Tbody>
        </Table>
    </Box>;
}
export default Properties;