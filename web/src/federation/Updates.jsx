import {
    Button, ButtonGroup, Center, CircularProgress, Flex, Menu,
    MenuButton, MenuItem, MenuList, Switch, Table, Tbody, Td, Text, Th, Thead, Tr,
    useToast,
    VStack
} from '@chakra-ui/react';
import { React, useContext, useEffect, useState } from "react";
import { BiChevronDown } from 'react-icons/bi';
import T from '../core/T';
import Server from '../server';
import UserContext from '../UserContext';
import Utils from '../core/utils';


function Updates(props) {
    const { project } = useContext(UserContext)
    const [range, setRange] = useState(1)
    const [state, setState] = useState(null)
    const { onClose } = props

    function getState() {
        const time = new Date()
        time.setDate(time.getDate() - range)
        Server.getFedState(project, JSON.stringify(time))
            .then(setState)
    }
    useEffect(getState, [range])

    function Row(props) {
        const { item } = props
        const [update, setUpdate] = useState(['new', 'newer'].includes(item.state))
        item.update = update

        function switchUpdate() {
            item.update = !update
            setUpdate(!update)
        }

        return <Tr>
            <Td>{item.path}</Td>
            <Td>{Utils.getFriendlyDate(item.modTime)}</Td>
            <Td>{item.user}</Td>
        </Tr>
    }


    const updates = state && state.updates && state.updates
        .sort((a, b) => (a.path.localeCompare(b.path)))
        .map(s => <Row key={s.path} item={s} />) || []
    const sent = state && state.sent && state.sent
        .sort((a, b) => (a.path.localeCompare(b.path)))
        .map(s => <Row key={s.path} item={s} />) || []

    const rangeButton = <Menu>
        <MenuButton colorScheme="blue" as={Button} rightIcon={<BiChevronDown />}>
            Range
        </MenuButton>
        <MenuList>
            <MenuItem onClick={_ => setRange(1)}>
                Today
            </MenuItem>
            <MenuItem onClick={_ => setRange(7)}>
                One Week
            </MenuItem>
            <MenuItem onClick={_ => setRange(30)}>
                30 Days
            </MenuItem>
            <MenuItem onClick={_ => setRange(356)}>
                One Year
            </MenuItem>
        </MenuList>
    </Menu>


    return state != null ? <VStack>
        <Text fontSize="xl" as="b">Received</Text>
        <Flex overflow="auto" w="100%">
            <Table overflow="auto" >
                <Thead>
                    <Tr>
                        <Th>File</Th>
                        <Th>Changed On</Th>
                        <Th>By</Th>
                    </Tr>
                </Thead>
                <Tbody>
                    {updates}
                </Tbody>
            </Table>
        </Flex>
        <Text fontSize="xl" as="b">Sent</Text>
        <Flex overflow="auto" w="100%">
            <Table overflow="auto" >
                <Thead>
                    <Tr>
                        <Th>File</Th>
                        <Th>Changed On</Th>
                        <Th>By</Th>
                    </Tr>
                </Thead>
                <Tbody>
                    {sent}
                </Tbody>
            </Table>
        </Flex>
        <ButtonGroup>
            {rangeButton}
            <Button onClick={onClose}><T>close</T></Button>
        </ButtonGroup>

    </VStack> :
        <Center><CircularProgress isIndeterminate color="green.300" size="100px" /></Center>
}
export default Updates;