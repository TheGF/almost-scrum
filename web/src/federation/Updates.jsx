import {
    Button, ButtonGroup, Center, CircularProgress, Flex, Menu,
    MenuButton, MenuItem, MenuList, Switch, Table, Tbody, Td, Th, Thead, Tr,
    useToast,
    VStack
} from '@chakra-ui/react';
import { React, useContext, useEffect, useState } from "react";
import { BiChevronDown } from 'react-icons/bi';
import T from '../core/T';
import Server from '../server';
import UserContext from '../UserContext';


function Updates(props) {
    const { project, reload } = useContext(UserContext)
    const [updates, setUpdates] = useState(null)
    const [importing, setImporting] = useState(false)
    const { onClose, exportSince } = props
    const toast = useToast()

    function getStatus() {
        Server.getFedState(project, true)
            .then(status => setUpdates(status.updates || []))
    }
    useEffect(getStatus, [])

    function importFiles() {
        setImporting(true)

        const filtered = updates.filter(u=>u.update)
                                .map(u=> ({...u, update: undefined}))

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



    const rows = []
    // updates && updates.filter(update => update.state != 'older')
    //     .sort((a, b) => a.loc.localeCompare(b.loc))
    //     .map(update => <Row key={update.loc} item={update} />)


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


    return updates != null ? <VStack>
        {rows.length ? <Flex overflow="auto" w="100%">
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
                    {rows}
                </Tbody>
            </Table>
        </Flex> : <h1>No updates</h1>}
        <ButtonGroup>
            <Button colorScheme="blue" onClick={importFiles} isLoading={importing}
                isDisabled={!rows.length}><T>import</T></Button>
            {exportSinceButton}
            <Button onClick={onClose}><T>close</T></Button>

        </ButtonGroup>

    </VStack> :
        <Center><CircularProgress isIndeterminate color="green.300" size="100px" /></Center>
}
export default Updates;