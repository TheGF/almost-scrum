import { Input, Select, Switch, Table, TableCaption, Tbody, Td, Tr, VStack } from '@chakra-ui/react';
import { React, useState, useContext } from "react";
import T from '../core/T';
import UserContext from '../UserContext';

function Properties(props) {
    const { task, saveTask } = props;
    const { info } = useContext(UserContext)
    const { property_model } = info
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
            return <Input value={value} onChange={onChange} />
        }

        function renderEnum() {
            const options = values && values.map(option =>
                <option value={option} key={option}>{T.translate(option)}</option>
            ) || []
            return <Select placeholder="Choose" value={value}
                size="small" onChange={onChange}>
                {options}
            </Select>
        }
        
        function renderTag() {
            const options = values && values.map(option =>
                <option value={option} key={option}>{T.translate(option)}</option>
            ) || []
            return <Select placeholder="Choose" value={value}
                size="small" onChange={onChange}>
                {options}
            </Select>
        }

        function renderUser() {
            const options = values && values.map(option =>
                <option value={option} key={option}>{T.translate(option)}</option>
            ) || []
            return <Select placeholder="Choose" value={value}
                size="small" onChange={onChange}>
                {options}
            </Select>
        }

        function renderBool() {
            return <Switch isChecked={value} onChange={onChange}/>    
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

    const rows = property_model.map(propertyDef => renderProperty(propertyDef))
    return <Table variant="striped" colorScheme="teal" size="sm">
        <TableCaption>Edit the properties of your project</TableCaption>
        <Tbody>
            {rows}
        </Tbody>
    </Table>;
}
export default Properties;