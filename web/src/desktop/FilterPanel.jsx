import { Button, HStack, Spacer } from "@chakra-ui/react";
import { React, useContext, useEffect, useRef, useState } from "react";
import { BsViewStacked, MdViewHeadline, RiFilterLine } from 'react-icons/all';
import ReactTags from 'react-tag-autocomplete';
import './reactTags.css'
import Server from '../server';
import UserContext from '../UserContext';

function FilterPanel(props) {
    const { project } = useContext(UserContext);
    const { compact, setCompact, setSearchKeys } = props;
    const [tags, setTags] = useState([])
    const [suggestions, setSuggestions] = useState([]);
    const reactTags = useRef()
    
    
    function updateSearchKeys() {
        const keys = tags.map(tag => tag.name)
        setSearchKeys(keys)
    }
    useEffect(updateSearchKeys, [tags])

    function getSuggestions(prefix) {
        Server.getSuggestions(project, prefix)
            .then(keys => setSuggestions(
                (keys || []).map((key, idx) => {
                    return { id: idx, name: key }
                })
            ))
    }

    function addTagToSearch(tag) {
        setTags([...tags, tag])
    }

    function deleteTagFromSearch(idx) {
        setTags(tags.filter((tag, i) => i != idx))
    }

    function onInputInSearch(query) {
        if (!query) return 
         
        if (query.startsWith('#') || query.startsWith('@') || query.length > 1) {
            getSuggestions(query)
        }
    }

    return <HStack spacing={3}>
        {/* <InputGroup>
            <InputLeftElement
                pointerEvents="none"
                children={<BsSearch color="gray.300" />}
            />

            <Input type="phone" placeholder="Search Filter" />
        </InputGroup> */}
        <ReactTags
            ref={reactTags}
            tags={tags}
            minQueryLength={1}
            suggestions={suggestions}
            onDelete={deleteTagFromSearch}
            onAddition={addTagToSearch}
            function onInput={onInputInSearch}
        />

        <Spacer />
        <Button><RiFilterLine /></Button>

        <Button onClick={_ => setCompact(true)} isActive={compact}>
            <MdViewHeadline />
        </Button>
        <Button onClick={_ => setCompact(false)} isActive={!compact}>
            <BsViewStacked />
        </Button>
    </HStack >
}

export default FilterPanel
