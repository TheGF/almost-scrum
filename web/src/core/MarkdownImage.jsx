

import { React, useState, useEffect } from "react";
import { Img } from '@chakra-ui/react';
import {
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalFooter,
    ModalBody,
    ModalCloseButton,
    useDisclosure,
} from "@chakra-ui/react"
import { Radio, RadioGroup, Button } from "@chakra-ui/react"
import { VStack, HStack, Stack, Box } from "@chakra-ui/react"
import {
    Slider,
    SliderTrack,
    SliderFilledTrack,
    SliderThumb,
} from "@chakra-ui/react"

function MarkdownImage(props) {
    const token = localStorage.token
    const { readOnly } = props
    const [alt, _align, _size] = parseOptions()
    const [align, setAlign] = useState(_align)
    const [size, setSize] = useState(_size)
    const { isOpen, onOpen, onClose } = useDisclosure()

    function parseOptions() {
        const options = props.alt && props.alt.split(' ') || []
        let alt2 = []
        let align = 'center'
        let size = 50

        for (const option of options) {
            if (['left', 'center', 'right'].includes(option)) {
                align = option
            }
            else if (option.endsWith('%')) {
                size = parseInt(option, 10)
            }
            else {
                alt2.push(option)
            }
        }
        return [alt2.join(' '), align, size]
    }

    function getMore() {
        let more = {
            style: { cursor: 'pointer' },
            htmlWidth: `${size}%`,
            htmlHeight: `${size}%`,
        }
        switch (align) {
            case 'right': more = { ...more, marginLeft: 'auto' }; break
            case 'center': more = { ...more, marginLeft: 'auto', marginRight: 'auto' }; break
        }
        return more
    }

    function save() {
        const orig = `[${props.alt}](${props.src})`
        const update = `[${alt} ${align} ${size}%](${props.src})`
        props.onUpdate && props.onUpdate(orig, update)
    }

    const onClick = readOnly ? null : onOpen
    const src = token ? `${props.src}?token=${token}` : props.src
    return <>
        <Img alt={alt}{...getMore()} src={src} onClick={onClick} />
        <Modal isOpen={isOpen} onClose={onClose}>
            <ModalOverlay />
            <ModalContent>
                <ModalHeader>Image Settings</ModalHeader>
                <ModalCloseButton />
                <ModalBody>
                    <VStack spacing={5}>
                        <HStack >
                            <RadioGroup value={align} onChange={setAlign} >
                                <Stack direction="row">
                                    <Radio value="left">Left</Radio>
                                    <Radio value="center">Center</Radio>
                                    <Radio value="right">Right</Radio>
                                </Stack>
                            </RadioGroup>
                        </HStack>
                        <HStack>
                            <label>Size</label>
                            <Slider min={5} max={100}
                                onChangeEnd={setSize}
                                defaultValue={size} w="200px">
                                <SliderTrack>
                                    <SliderFilledTrack />
                                </SliderTrack>
                                <SliderThumb />
                            </Slider>
                            <label>{size}%</label>
                        </HStack>
                    </VStack>
                </ModalBody>
                <ModalFooter>
                    <Button colorScheme="blue" mr={3} onClick={save}>
                        Save
                    </Button>
                </ModalFooter>
            </ModalContent>
        </Modal>
    </>
}

export default MarkdownImage