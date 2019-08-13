import React from 'react';
import xid from 'xid-js';
import './App.css';
import Result from "./Result";

const numImgs = 10;

class App extends React.PureComponent {
    constructor(props) {
        super(props);
        this.state = {
            id: null,
            running: false,
            items: [],
            results: null,
        };
    }

    render() {
        return (
            <div className="App">
                <h1>DNS Leak Test</h1>

                {this.state.running && !this.done() && <p>Running...</p>}

                <button
                    onClick={() => this.start()}
                    disabled={this.state.running}
                >
                    Start test
                </button>

                {this.state.running && this.state.items.map((item) => (
                    <img
                        alt=""
                        key={item.num}
                        style={{display: 'none'}}
                        src={item.src}
                        onError={this.onError(item.num)}
                    />
                ))}

               <div className="flex mt-4">
                   {this.done() && (
                       <table className="pure-table">
                           <thead>
                           <tr>
                               <th>Number</th>
                               <th>IP</th>
                               <th>Other detail</th>
                           </tr>
                           </thead>
                           <tbody>
                           {this.state.results && this.state.results.results.map((r) => {
                               return r.ips.map((ip) => (
                                   <Result key={ip} number={r.number} {...ip} />
                               ))
                           })}
                           </tbody>
                       </table>
                   )}
               </div>
            </div>
        );
    }

    done() {
        return this.state.items.filter((item) => !item.done).length === 0;
    }

    onError = (idx) => () => {
        const items = [...this.state.items];
        items[idx].done = true;
        this.setState({items}, () => {
            if (this.done()) {
                this.getResults();
            }
        });
    };

    getResults = () => {
        fetch(`https://dnsleak.app/api/results/456`)
            .then((r) => r.json())
            .then((data) => {
                this.setState({results: data});
            })
            .catch(console.error);
    };

    start() {
        const id = xid.next();
        const items = new Array(numImgs).fill(0).map((_, idx) => ({
            src: `https://${idx}.${id}.dnsleak.app/image.png`,
            num: idx,
            done: false,
        }));

        this.setState({
            id,
            items,
            running: true,
        })
    }
}

export default App;
