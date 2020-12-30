import { HStack, Input, InputGroup, InputLeftElement, Spacer, Flex, Button } from "@chakra-ui/react";
import { React, useEffect, useState, useContext } from "react";
import { BsSearch, BsViewStacked, MdViewHeadline } from 'react-icons/all';

function FilterPanel(props) {

    return <HStack spacing={3}>
        <InputGroup>
            <InputLeftElement
                pointerEvents="none"
                children={<BsSearch color="gray.300" />}
            />
            <Input type="phone" placeholder="Search Filter" />
        </InputGroup>
        <Spacer />
        <Button size="sm"><MdViewHeadline/></Button>
        <Button size="sm"><BsViewStacked /></Button>
    </HStack>
}

export default FilterPanel
