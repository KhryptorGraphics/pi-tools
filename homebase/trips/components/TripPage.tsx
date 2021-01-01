import React from "react";
import {Helmet} from "react-helmet";
import {useHistory, useParams} from "react-router-dom";
import {format, formatDuration, intervalToDuration, parseISO} from "date-fns";
import {DescriptionField} from "com_github_mjm_pi_tools/homebase/components/DescriptionField";
import {TripTagField} from "com_github_mjm_pi_tools/homebase/trips/components/TripTagField";
import {Alert} from "com_github_mjm_pi_tools/homebase/components/Alert";
import {graphql, useLazyLoadQuery} from "react-relay/hooks";
import {TripPageQuery} from "com_github_mjm_pi_tools/homebase/api/__generated__/TripPageQuery.graphql";
import {useIgnoreTrip} from "com_github_mjm_pi_tools/homebase/trips/lib/IgnoreTrip";

export function TripPage() {
    const {id} = useParams<{ id: string }>();
    const data = useLazyLoadQuery<TripPageQuery>(
        graphql`
            query TripPageQuery($id: ID!) {
                viewer {
                    trip(id: $id) {
                        id
                        leftAt
                        returnedAt
                        ...TripTagField_trip
                    }
                }
            }
        `,
        {id},
    );

    const trip = data.viewer.trip;

    return (
        <main className="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
            <Helmet>
                <title>Trip Details</title>
            </Helmet>

            <div className="bg-white shadow overflow-hidden sm:rounded-lg">
                <div className="bg-white px-4 py-5 border-b border-gray-200 sm:px-6">
                    <div className="-ml-4 -mt-2 flex items-center justify-between flex-wrap sm:flex-nowrap">
                        <div className="ml-4 mt-2">
                            <h3 className="text-lg leading-6 font-medium text-gray-900">
                                Trip Details
                            </h3>
                        </div>
                        <div className="ml-4 mt-2 flex-shrink-0 flex">
                            {trip && <TripIgnoreButton id={trip.id}/>}
                        </div>
                    </div>
                </div>
                {trip ? (
                    <div>
                        <dl>
                            <DescriptionField label="Left at" offset>
                                {format(parseISO(trip.leftAt), "PPpp")}
                            </DescriptionField>
                            {trip.returnedAt && (
                                <>
                                    <DescriptionField label="Returned at">
                                        {format(parseISO(trip.returnedAt), "PPpp")}
                                    </DescriptionField>
                                    <DescriptionField label="Duration" offset>
                                        {formatDuration(intervalToDuration({
                                            start: parseISO(trip.leftAt),
                                            end: parseISO(trip.returnedAt),
                                        }))}
                                    </DescriptionField>
                                </>
                            )}
                            <DescriptionField label="Tags">
                                <TripTagField trip={trip}/>
                            </DescriptionField>
                        </dl>
                    </div>
                ) : (
                    <Alert title="Couldn't load this trip" severity="error" rounded={false}>
                        No trip was found with this ID.
                    </Alert>
                )}
            </div>
        </main>
    );
}

function TripIgnoreButton({id}: { id: string }) {
    const history = useHistory();
    const [commit, isInFlight] = useIgnoreTrip();

    async function onIgnore() {
        try {
            await commit(id);

            // return to the trips page upon successful ignore
            history.push("/trips");
        } catch (e) {
            console.error(e);
        }
    }

    return (
        <span className="inline-flex rounded-md shadow-sm">
            <button type="button"
                    disabled={isInFlight}
                    onClick={onIgnore}
                    className="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm leading-5 font-medium rounded-md text-gray-700 bg-white hover:text-gray-500 focus:outline-none focus:ring-blue focus:border-blue-300 active:bg-gray-50 active:text-gray-800">
                <svg className="-ml-1 mr-2 h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg"
                     viewBox="0 0 20 20"
                     fill="currentColor">
                    <path fillRule="evenodd"
                          d="M10 18a8 8 0 100-16 8 8 0 000 16zM7 9a1 1 0 000 2h6a1 1 0 100-2H7z"
                          clipRule="evenodd"/>
                </svg>
                <span>Ignore</span>
            </button>
        </span>
    );
}
