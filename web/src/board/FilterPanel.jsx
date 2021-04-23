import {
    Button, HStack, Popover,

    PopoverContent, PopoverTrigger, Spacer
} from "@chakra-ui/react";
import { React, useContext, useEffect, useRef, useState } from "react";
import { BsViewStacked, MdViewHeadline, RiFilterLine } from 'react-icons/all';
import ReactTags from 'react-tag-autocomplete';
import Server from '../server';
import UserContext from '../UserContext';
import Filter from "./Filter";
import NewTask from "./NewTask";
import './reactTags.css';

function FilterPanel(props) {
    const { project, info } = useContext(UserContext);
    const { compact, toggleCompact, setSearchKeys, onNewTask, users, board } = props;
    const [tags, setTags] = useState([])
    const [suggestions, setSuggestions] = useState([]);
    const [showFilter, setShowFilter] = useState(false);
    const reactTags = useRef()


    function updateSearchKeys() {
        const keys = tags.map(tag => tag.name)
        setSearchKeys(keys)
    }
    useEffect(updateSearchKeys, [tags])

    function getSuggestions(prefix) {
        Server.getSuggestions(project, prefix)
            .then(keys => setSuggestions(
                (keys || []).map(key => {
                    return { id: key, name: key }
                })
            ))
    }

    function addTagToSearch(tag) {
        setTags([...tags, tag])
    }

    function deleteTagFromSearch(id) {
        setTags(tags.filter((tag, i) => i != id))
    }

    function onInputInSearch(query) {
        if (!query) return

        if (query.startsWith('#') || query.startsWith('@') || query.length > 1) {
            getSuggestions(query)
        }
    }

    const compactButton = compact ?
        <Button onClick={toggleCompact} isActive={!compact}
            title="Show all task content">
            <BsViewStacked />
        </Button> :
        <Button onClick={toggleCompact} isActive={compact}
            title="Show only tasks header">
            <MdViewHeadline />
        </Button>


    const filterButton = <Popover placement="bottom-end" maxW={500} className="filter-pop" variant="responsive">
        <PopoverTrigger>
            <Button onClick={_ => setShowFilter(!showFilter)}>
                <RiFilterLine />
            </Button>
        </PopoverTrigger>
        <PopoverContent maxW={500}>
            <Filter users={users} tags={tags} setTags={setTags}/>
        </PopoverContent>
    </Popover>

    return <HStack spacing={3}>
        <NewTask board={board} onNewTask={onNewTask}/>
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
        {compactButton}
        {filterButton}
    </HStack >

}

export default FilterPanel
