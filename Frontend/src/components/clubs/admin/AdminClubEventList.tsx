import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import api from "../../../utils/api";

interface Events {
    id: string;
    name: string;
    date: string;
    description: string;
    begin_time: string;
    end_time: string;
}

const AdminClubEventList = () => {

    const { id } = useParams();
    const [events, setEvents] = useState<Events[]>([]);

    useEffect(() => {
        const fetchEvents = async () => {
            try {
                const response = await api.get(`/api/v1/clubs/${id}/events`);
                setEvents(response.data);
            } catch (error) {
                console.error("Error fetching events:", error);
                setEvents([]);
            }
        };
        fetchEvents();
    }, [id]);

    return (
        <div>
            <h3>Events</h3>
            <table className="basic-table">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Date</th>
                        <th>Description</th>
                        <th>Begin Time</th>
                        <th>End Time</th>
                    </tr>
                </thead>
                <tbody>
                    {events.map((event) => (
                        <tr key={event.id}>
                            <td>{event.name}</td>
                            <td>{event.date}</td>
                            <td>{event.description}</td>
                            <td>{event.begin_time}</td>
                            <td>{event.end_time}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}

export default AdminClubEventList;