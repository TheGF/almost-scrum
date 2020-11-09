const emptyStory = {
    description: '',
    points: 0,
    status: 'Draft',
    owner: '',
    tasks: [],
    timeline: [],
    attachments: [],
}

const pointsUnitOptions = {
    fibo: [1, 2, 3, 5, 8, 13, 21, 34],
    linear: [1, 2, 3, 4, 5, 6, 7, 8, 9]
}

const statusOptions = ['Draft', 'Ready', 'Build', 'Done'];

export { emptyStory, pointsUnitOptions, statusOptions } ;