import {
    Button, HStack, Img, Input, Modal, ModalBody, ModalCloseButton, ModalContent,
    ModalFooter, ModalHeader, ModalOverlay, Radio, RadioGroup, Slider,
    SliderFilledTrack, SliderThumb, SliderTrack, Stack, useDisclosure, VStack
} from '@chakra-ui/react';
import { React, useContext, useState } from "react";
import UserContext from '../UserContext';


function parseUrl(url) {
    const sharp = url.lastIndexOf('#');
    if (sharp === -1) {
        return [url, 'left', 100, null]
    }
    const options = url.substring(sharp + 1).split(',');
    let align = 'left'
    let size = 50
    let caption = null

    for (const option of options) {
        const parts = option.split('=')
        if (parts.length !== 2) continue

        switch (parts[0]) {
            case 'align': align = parts[1]
            case 'size': if (isNaN(parts[1]) == false) {
                size = parseInt(parts[1], 10)
            }
            case 'caption': caption = decodeURIComponent(parts[1])
        }
    }
    return [url.substring(0, sharp), align, size, caption]
}

const alignToStyle = {
    left: 'cursor:pointer',
    center: 'cursor:pointer;margin-left:auto;margin-right:auto',
    right: 'cursor:pointer;margin-left:auto',
}

function getImg(url, id) {
    const [_, align, size, caption] = parseUrl(url)

    let style = alignToStyle[align]
    return `<img id="${id}" src="${url}" style="${style}" \
            width="${size}%" title="${caption}"/>`
}


function MarkdownImage(props) {
    const { project } = useContext(UserContext);
    const { readOnly, image, setImage } = props
    if (!image) {
        return null
    }

    const [url_, align_, size_, caption_] = parseUrl(props.image.src)
    const [align, setAlign] = useState(align_);
    const [size, setSize] = useState(size_);
    const [caption, setCaption] = useState(caption_);

    image.src = `${url_}#size=${size},align=${align},caption=${encodeURIComponent(caption)}`
    image.style = alignToStyle[align]
    image.width = image.parentElement.clientWidth * size / 100
    image.title = caption

    // function update() {
    //     const [url, _, __, ___] = parseUrl(image.src)

    //     image.src = `${url}#size=${size},align=${align}`
    //     image.style = alignToStyle[align]
    //     image.width = image.parentElement.clientWidth * size / 100
    // }

    const onClick = readOnly ? null : e => {
        e.stopPropagation()
        onOpen(true)
    }



    return <>
        <Modal isOpen={image} onClose={_ => setImage(null)}>
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
                        <HStack>
                            <label>Caption</label>
                            <Input value={caption}
                                onChange={e => setCaption(e.target.value)} />
                        </HStack>
                    </VStack>
                </ModalBody>
                <ModalFooter>
                    {<Button colorScheme="blue" mr={3} onClick={_ => setImage(null)}>
                        Close
                    </Button>}
                </ModalFooter>
            </ModalContent>
        </Modal>
    </>
}

MarkdownImage.parseUrl = parseUrl
MarkdownImage.getImg = getImg

export default MarkdownImage