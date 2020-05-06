import React, { Component } from 'react';
import { Spinner, Table, FormattedDate, RunDropdown } from '@tektoncd/dashboard-components';
import { Link } from 'react-router-dom';
import { urls } from '@tektoncd/dashboard-utils';
import { IntlProvider } from "react-intl";
import { getProwjobs } from './api/index';
import {
  CheckmarkFilled20 as CheckmarkFilled,
  CloseFilled20 as CloseFilled,
  Time20 as Time
} from '@carbon/icons-react';

import './status.scss'

class WebhooksApp extends Component {
  state = {
    prowjobs: null,
    loading: false
  };

  componentDidMount() {
    this.loadProwjobs();
  };

  async loadProwjobs() {
    this.setState({ loading: true });
    try {
      const prowjobs = await getProwjobs();
      this.setState({ prowjobs });
    } catch (error) { }
    this.setState({ loading: false });
  };

  delete = prowjob => {
    console.log(prowjob);
  };

  viewYaml = prowjob => {
    console.log(prowjob);
  };

  viewOnProw = prowjob => {
    console.log(prowjob);
  };

  render() {
    const { prowjobs, loading } = this.state;

    const headers = [
      { key: 'status', header: 'Status' },
      { key: 'name', header: 'Name' },
      { key: 'createdTime', header: 'Created' },
      { key: 'type', header: 'Type' },
      { key: 'repo', header: 'Repository' },
      { key: 'agent', header: 'Agent' },
      { key: 'actions', header: '' }
    ];

    const actions = [
      { actionText: "Delete", action: this.delete },
      { actionText: "View Yaml", action: this.viewYaml },
      { actionText: "View on Prow", action: this.viewOnProw }
    ];

    const createPipelineRunURL = prowjob => urls.pipelineRuns.byName({
      namespace: prowjob.metadata.namespace,
      pipelineRunName: prowjob.metadata.name
    });

    const getProwjobStatusIcon = (status, reason) => {
      if (!status || (status === 'pending')) {
        return <Time className="status-icon" />;
      }

      if (status === 'running') {
        return <Spinner className="status-icon" />;
      }
    
      let Icon;
      if (status === 'success') {
        Icon = CheckmarkFilled;
      } else if (status === 'failure') {
        Icon = CloseFilled;
      }

      return Icon ? <Icon className="status-icon" /> : null;
    };

    const getProwjobStatus = prowjob => {
      const reason = prowjob.status.description;
      const status = prowjob.status.state;

      return (
        <div className="definition">
          <div
            className="status"
            data-status={status}
            data-reason={reason}
            title={reason}
          >
            {getProwjobStatusIcon(status, reason)}
          </div>
        </div>
      );
    };

    const rows = [];

    if (prowjobs) {
      prowjobs.items.forEach(prowjob => {
        rows.push({
          id: prowjob.metadata.name,
          status: getProwjobStatus(prowjob),
          name:<Link to={createPipelineRunURL(prowjob)} title={prowjob.metadata.name}>{prowjob.metadata.name}</Link>,
          type: prowjob.spec.type,
          repo: `${prowjob.spec.refs.org}/${prowjob.spec.refs.repo}`,
          agent: prowjob.spec.agent,
          createdTime: <FormattedDate date={prowjob.metadata.creationTimestamp} relative />,
          actions:  <RunDropdown items={actions} resource={prowjob} />
        });
      });
    }

    return (
      <IntlProvider locale="en">
        <h1>Prowjobs</h1>
        <br/>
        <Table
          headers={headers}
          rows={rows}
          loading={loading}
          emptyTextAllNamespaces='No prow jobs'
          emptyTextSelectedNamespace='No prow jobs'
        />
      </IntlProvider>
    );
  };
}

export default WebhooksApp;