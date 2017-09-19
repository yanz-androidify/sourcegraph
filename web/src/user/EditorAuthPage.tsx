import * as React from 'react'
import { PageTitle } from 'sourcegraph/components/PageTitle'
import { sourcegraphContext } from 'sourcegraph/util/sourcegraphContext'

/**
 * Page to enable users to authenticate/link to their editors
 */
export class EditorAuthPage extends React.Component {
    public render(): JSX.Element | null {
        return (
            <div>
                <PageTitle title='authenticate editor' />
                <div>Welcome to Sourcegraph!</div>
                <p>Your session ID is: <span>{sourcegraphContext.xhrHeaders.Authorization}</span></p>
            </div>
        )
    }
}
