/**
 * Utility functions for event-related operations
 */

export interface RSVPCounts {
    yes?: number;
    no?: number;
    maybe?: number;
}

export interface EventRSVP {
    Response: string;
}

/**
 * Calculates RSVP counts by response type from a list of RSVPs.
 * Validates that responses are one of the allowed values ('yes', 'no', 'maybe').
 * 
 * @param rsvpList - Array of RSVP objects with Response field
 * @returns Object with counts for each valid response type
 */
export function calculateRSVPCounts(rsvpList: EventRSVP[]): RSVPCounts {
    const allowedResponses: (keyof RSVPCounts)[] = ['yes', 'no', 'maybe'];
    
    return rsvpList.reduce((acc: RSVPCounts, rsvp: EventRSVP) => {
        const responseKey = rsvp.Response.toLowerCase() as keyof RSVPCounts;
        
        // Only count valid response types
        if (allowedResponses.includes(responseKey)) {
            acc[responseKey] = (acc[responseKey] || 0) + 1;
        }
        
        return acc;
    }, {});
}
