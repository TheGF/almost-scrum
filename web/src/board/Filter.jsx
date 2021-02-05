import { Box, Checkbox, CheckboxGroup, FormControl, FormLabel, HStack, ButtonGroup, Button, Text } from "@chakra-ui/react";
import { React, useContext } from "react";
import UserContext from '../UserContext';
import './reactTags.css';

function Filter(props) {
    const { info } = useContext(UserContext);
    const { propertyModel } = info
    const { users, tags, setTags } = props

    function switchTag(value) {
        const filtered = tags.filter(t => t.id != value)
        if (filtered.length < tags.length) {
            setTags(filtered)
        } else {
            setTags([...tags, { id: value, name: value }])
        }
    }

    function getTagFilter(name, values) {
        const items = values.map(v =>
            <Button isActive={tags.filter(t => t.id == v).length}
                onClick={_ => switchTag(v)}>{v}</Button>
        )
        return <HStack>
            <Text>{name}</Text>
            <ButtonGroup colorScheme="green" size="sm" >
                {items}
            </ButtonGroup>
        </HStack>
    }

    function getFieldFilter(field) {
        switch (field.kind) {
            case 'User': return getTagFilter(field.name, users)
            case 'Tag': return getTagFilter(field.name, field.values)
            default: return null
        }
    }

    const fields = propertyModel.filter(model => model.isFilter)
        .map(getFieldFilter)

    return <Box margin={3}>
        {fields}
    </Box>
}

export default Filter
