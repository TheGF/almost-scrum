import {
    Button, ButtonGroup, Center, CircularProgress, Flex, Menu,
    MenuButton, MenuItem, MenuList, Switch, Table, Tbody, Td, Text, Th, Thead, Tr,
    useToast, VStack
} from '@chakra-ui/react';
import { React, useContext, useEffect, useState } from "react";
import { BiChevronDown } from 'react-icons/bi';
import T from '../core/T';
import Server from '../server';
import UserContext from '../UserContext';


function Parked(props) {
    const { project, reload } = useContext(UserContext)
    const [state, setState] = useState(null)
    const [importing, setImporting] = useState(false)
    const { onClose, exportSince } = props
    const [update, setUpdate] = useState([])
    const [ignore, setIgnore] = useState([])
    const toast = useToast()

    function getState() {
        const time = new Date()
        Server.getFedState(project, JSON.stringify(time))
            .then(setState)
    }
    useEffect(getState, [])

    function importFiles() {
        setImporting(true)

        const filtered = updates.filter(u => u.update)
            .map(u => ({ ...u, update: undefined }))

        if (filtered.length) {
            Server.postFedImport(project, filtered)
                .then(_ => setImporting(false))
                .then(reload)
        }
    }

    function Row(props) {
        const { item } = props
        const [update, setUpdate] = useState(['new', 'newer'].includes(item.state))
        item.update = update

        function switchUpdate() {
            item.update = !update
            setUpdate(!update)
        }

        return <Tr>
            <Td>{item.loc}</Td>
            <Td>{item.owner}</Td>
            <Td>{item.state}</Td>
            <Td><Switch isChecked={update} onChange={switchUpdate} /></Td>
        </Tr>
    }



    const parked = state && state.parked && state.parked
        .sort((a, b) => (a.path.localeCompare(b.path)))
        .map(s => <Row key={s.path} item={s} />)
    // updates && updates.filter(update => update.state != 'older')
    //     .sort((a, b) => a.loc.localeCompare(b.loc))
    //     .map(update =>)


    const exportSinceButton = <Menu>
        <MenuButton colorScheme="blue" as={Button} rightIcon={<BiChevronDown />}>
            Export Since
        </MenuButton>
        <MenuList>
            <MenuItem onClick={_ => exportSince('today')}>
                Today
            </MenuItem>
            <MenuItem onClick={_ => exportSince('week')}>
                One Week
            </MenuItem>
            <MenuItem onClick={_ => exportSince('month')}>
                One Month
            </MenuItem>
            <MenuItem onClick={_ => exportSince('all')}>
                The Big Bang
            </MenuItem>
        </MenuList>
    </Menu>


    return state != null ? <VStack>
        <Text fontSize="xl" as="b">Parked</Text>
        <Flex overflow="auto" w="100%">
            <Table overflow="auto" >
                <Thead>
                    <Tr>
                        <Th>File</Th>
                        <Th>From</Th>
                        <Th>State</Th>
                        <Th>Update</Th>
                    </Tr>
                </Thead>
                <Tbody>
                    {parked}
                </Tbody>
            </Table>
        </Flex>
        <ButtonGroup>
            <Button colorScheme="blue" onClick={importFiles} isLoading={importing}
                isDisabled={!parked}><T>import</T></Button>
            {exportSinceButton}
            <Button onClick={onClose}><T>close</T></Button>
        </ButtonGroup>
    </VStack> :
        <Center><CircularProgress isIndeterminate color="green.300" size="100px" /></Center>
}
export default Parked;