import { Button, HStack, Input, InputGroup, InputLeftElement, Spacer } from "@chakra-ui/react";
import { React } from "react";
import { BsSearch, BsViewStacked, MdViewHeadline, RiFilterLine } from 'react-icons/all';

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
        <Button size="sm"><RiFilterLine /></Button>

        <Button size="sm" onClick={_ => setCompact(true)} isActive={compact}>
            <MdViewHeadline />
        </Button>
        <Button size="sm" onClick={_ => setCompact(false)} isActive={!compact}>
            <BsViewStacked />
        </Button>
    </HStack >
}

export default FilterPanel
