import { HStack, Input, InputGroup, InputLeftElement, Spacer, Flex, Button } from "@chakra-ui/react";
import { React, useEffect, useState, useContext } from "react";
import { BsSearch, BsViewStacked, MdViewHeadline } from 'react-icons/all';

function FilterPanel(props) {
    const {compact, setCompact} = props;

    return <HStack spacing={3}>
        <InputGroup>
            <InputLeftElement
                pointerEvents="none"
                children={<BsSearch color="gray.300" />}
            />
            <Input type="phone" placeholder="Search Filter" />
        </InputGroup>
        <Spacer />
        <Button size="sm" onClick={_ => setCompact(true)} isActive={compact}>
            <MdViewHeadline />
        </Button>
        <Button size="sm" onClick={_ => setCompact(false)} isActive={!compact}>
            <BsViewStacked />
        </Button>
    </HStack >
}

export default FilterPanel
