import {
    Box, FormLabel, HStack, Input, Select,
    Slider, SliderFilledTrack, SliderThumb, SliderTrack, Spacer, Switch, Table, Tbody,
    Td, Tr
} from '@chakra-ui/react';
import { React, useContext, useState } from "react";
import ReactDatePicker from 'react-datepicker';
import T from '../core/T';
import UserContext from '../UserContext';
import TasksSelector from './TasksSelector';


function Properties(props) {
    const { task, saveTask, readOnly, users, height } = props;
    const { info } = useContext(UserContext)
    const { properties } = task;
    const type = properties && properties['Type'] || null

    const model = info && info.models &&
        info.models.filter(m => m.name == type).shift() || []

    function renderProperty(propertyDef) {
        const { name, kind, values } = propertyDef
        const [value, setValue] = useState(properties[name] || '');

        function onChange(evt) {
            const value = evt && evt.target && evt.target.value;
            if (value != undefined) changeValue(value)
        }

        function changeValue(v) {
            v = `${v}`
            if (v != value) {
                properties[name] = v
                setValue(v)
                saveTask(task)
            }
        }

        function renderString() {
            return <Input readOnly={readOnly} value={value}
                onChange={onChange} size="xs" />
        }

        function renderPercentage() {
            let p = parseInt(value, 10)
            if ((p >= 0 && p <= 100) == false) {
                p = 0
            }

            return <HStack spacing={2}>
                <label>{p}%</label>
                {readOnly ? null :
                    <Slider min={0} max={100} onChangeEnd={changeValue}
                        defaultValue={p} maxW="80%" >
                        <SliderTrack>
                            <SliderFilledTrack />
                        </SliderTrack>
                        <SliderThumb />
                    </Slider>
                }
            </HStack>
        }

        function renderTasks(maxSize) {
            return <TasksSelector value={value} onChange={changeValue}
                maxSize={maxSize} />
        }

        function renderDate() {
            const startDate = Date.parse(value)
            return <ReactDatePicker selected={startDate} readOnly={readOnly}
                onChange={date => changeValue(date && date.toISOString())} />
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
            case 'Date': input = renderDate(); break
            case 'Percentage': input = renderPercentage(); break
            case 'Tasks': input = renderTasks(16); break
            case 'Task': input = renderTasks(1); break
        }

        return <Tr key={name}>
            <Td><T>{name}</T></Td>
            <Td>{input}</Td>
        </Tr>
    }

    const rows = model.properties.map(propertyDef => renderProperty(propertyDef))
    if (type == null) {
        return <font size="md">Corrupted task: no Type found</font>
    }
    if (model.length == 0) {
        return <font size="md">Corrupted task: invalid Type {type}</font>
    }

    return <Box maxH={height - 50} style={{ overflowY: 'auto' }}>
        <Table size="sm">
            <Tbody>
                {rows}
            </Tbody>
        </Table>
    </Box>
}
export default Properties;