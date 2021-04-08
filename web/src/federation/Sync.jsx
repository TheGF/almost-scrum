import {
    Button, ButtonGroup, Center, CircularProgress, Flex, Switch, Table, Tbody, Td, Th, Thead, Tr,
    VStack
} from '@chakra-ui/react';
import { React, useContext, useEffect, useState } from "react";
import T from '../core/T';
import Server from '../server';
import UserContext from '../UserContext';


function Sync(props) {
    const { project, reload } = useContext(UserContext)
    const [diffs, setDiffs] = useState(null)
//    const [exported, setExported] = useState(null)
    const [updated, setUpdated] = useState(false)
    const { onClose } = props

    function getDiffs() {
        Server.getFedDiffs(project, true)
            .then(setDiffs)
    }
    useEffect(getDiffs, [])

    function importFiles() {
        setDiffs(null)
        Server.postFedImport(project, diffs)
            .then(_ => setUpdated(true))
            .then(reload)
    }

    // function exportFiles() {
    //     Server.postFedExport(project, diffs)
    //         .then(setExported)
    // }


    function Row(props) {
        const { name, item, header } = props
        const [strategy, setStrategy] = useState(item.strategy)

        function switchStrategy() {
            item.strategy = strategy === 'extract' ? 'ignore' : 'extract'
            setStrategy(item.strategy)
        }

        return <Tr>
            <Td>{name}</Td>
            <Td>{header.user}@{header.hostname}</Td>
            <Td>{item.match}</Td>
            <Td><Switch isChecked={strategy == 'extract'} onChange={switchStrategy} /></Td>
        </Tr>
    }

    const rows = diffs && diffs.flatMap(log => Object.keys(log.items).sort()
        .filter(name => log.items[name].match != 'outdated')
        .map(name => <Row key={name} name={name} item={log.items[name]} header={log.header} />))

    return diffs != null ? <VStack>
        {rows.length ? <Flex overflow="auto" w="100%">
            <Table overflow="auto" >
                <Thead>
                    <Tr>
                        <Th>File</Th>
                        <Th>From</Th>
                        <Th>Match</Th>
                        <Th>Update</Th>
                    </Tr>
                </Thead>
                <Tbody>
                    {rows}
                </Tbody>
            </Table>
        </Flex> : <h1>No updates</h1>}
        {/* {exported && exported.length ? <h1>
            Exported {exported.length} file{exported.length > 1 ? 's' : ''}: {exported.join(',')}
        </h1> : <h1>No exports</h1>} */}
        <ButtonGroup>
            <Button colorScheme="blue" onClick={importFiles} isDisabled={!rows.length}><T>import</T></Button>
            <Button onClick={onClose}><T>close</T></Button>
        </ButtonGroup>

    </VStack> :
        <Center><CircularProgress isIndeterminate color="green.300" size="100px" /></Center>
}
export default Sync;