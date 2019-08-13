import React from 'react';

const Result = ({ number, address, info }) => {
    return (
        <tr>
            <td>{number}</td>
            <td>{address}</td>
            <td>{info ? "hi" : "n/a"}</td>
        </tr>
    );
};

export default Result;
