import {
    Box, Button, ButtonGroup, FormControl, FormLabel, Image, Textarea, useToast, VStack, Wrap, WrapItem
} from "@chakra-ui/react";
import { React, useContext, useEffect, useState } from "react";
import Server from "../server";

const images = ['banana', 'beer', 'burger', 'chocolate', 'coke', 'cornflakes',
    'cupcake', 'egg', 'fries', 'hotdog', 'juice', 'muffin', 'orange',
    'pasta', 'pizza', 'popcorn', 'steak', 'sushi', 'water', 'wine'
]


function Join(props) {
    const { onClose, project } = props
    const [token, setToken] = useState(props.token || null)
    const [selected, setSelected] = useState([])
    const toast = useToast()

    function lookForInvite() {
        const queryParams = new URLSearchParams(window.location.search);
        const invite = queryParams.get('invite');
        const token = queryParams.get('token');

        if (invite == "" && token) {
            setToken(token)
            onOpen()
        }
    }
    useEffect(lookForInvite, [])

    function join() {
        const key = selected.join(',')
        Server.postFedJoin(project, key, token).then(_ => {
            toast({
                title: `Claim Success`,
                description: 'The invite has been successfully claimed',
                status: "success",
                duration: 9000,
                isClosable: true,
            })
        })

        onClose()
    }

    function clickIcon(name) {
        if (selected.includes(name)) {
            setSelected(selected.filter(n => n != name))
        } else {
            setSelected([...selected, name])
        }
    }

    const icons = images.map(n => <WrapItem>
        <Box bg={selected.includes(n) ? 'yellow' : null} borderWidth={1} p={1}
            onClick={_ => clickIcon(n)}>
            <VStack>
                <Image boxSize="50px" src={`/icons/${n}.svg`} ></Image>
                <label>{n}</label>
            </VStack>
        </Box>
    </WrapItem>)

    return <VStack>
        <FormControl isRequired>
            <FormLabel>What had the Gopher for dinner?</FormLabel>
            <Wrap>
                {icons}
            </Wrap>

        </FormControl>
        <FormControl isRequired>
            <FormLabel>Token</FormLabel>
            <Textarea rows={12} value={token} onChange={e => setToken(e.target.value)} />
        </FormControl>
        <ButtonGroup>
            <Button colorScheme="blue" isDisabled={token == null || token.length == 0 || selected.length != 2}
                onClick={join}>Confirm</Button>
            <Button onClick={onClose}>Close</Button>
        </ButtonGroup>
    </VStack>
}

export default Join