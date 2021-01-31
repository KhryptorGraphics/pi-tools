import SwiftUI
import RelaySwiftUI
import detect_presence_ios_relay_generated

private let query = graphql("""
query RecordingTabQuery {
    # workaround for the relay compiler, it needs some kind of server field present
    # in the query.
    ...on Query { __typename }

    appEvents {
        id
        ...AppEventRow_event
    }
}
""")

struct RecordingTab: View {
    @ObservedObject var model: AppModel
    @Query<RecordingTabQuery>(fetchPolicy: .storeOnly) var query
    @AppStorage("recordToDevServer") private var recordToDevServer: Bool = false

    var body: some View {
        List {
            Section {
                Toggle("Record to development server", isOn: $recordToDevServer)

                if recordToDevServer {
                    Button {
                        model.beginTrip()
                    } label: {
                        Label("Simulate begin trip", systemImage: "play.fill")
                    }

                    Button {
                        model.endTrip()
                    } label: {
                        Label("Simulate end trip", systemImage: "stop.fill")
                    }
                }

                if let trip = model.currentTrip {
                    Text("Current trip started ") + Text(trip.leftAt, style: .relative) + Text(" ago")
                } else {
                    Text("Not currently on a trip")
                }

                if model.queuedTripCount > 0 {
                    Button {
                        model.recordQueuedTrips()
                    } label: {
                        Label("Record \(model.queuedTripCount) queued trips", systemImage: "icloud.and.arrow.up.fill")
                    }

                    Button {
                        model.clearQueuedTrips()
                    } label: {
                        Label("Clear queued trips", systemImage: "trash.fill")
                    }
                }
            }

            switch query.get() {
            case .loading:
                Text("Loading (this shouldn't happen)")
            case .failure(let error):
                Text("Error: \(error.localizedDescription)")
            case .success(let data):
                if let data = data {
                    Section(header: Text("All Events")) {
                        ForEach(data.appEvents ?? []) { event in
                            AppEventRow(event: event.asFragment())
                        }
                    }
                } else {
                    Text("No data")
                }
            }
        }
        .navigationTitle("Presence")
    }
}
